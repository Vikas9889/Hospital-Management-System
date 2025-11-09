# Seed demo data from hms_seed_data CSVs into running services
# Usage: from repo root PowerShell: .\scripts\seed-demo.ps1
# This script:
#  - posts patients to user-service (/v1/users)
#  - posts doctors to doctor-service (/v1/doctors)
#  - posts appointments to appointment-service (/v1/appointments) mapping CSV ids to created ids
#  - for appointments with status COMPLETED, it calls the complete endpoint
#  - inserts prescriptions into the appointments DB (prescriptions table)
#  - creates bills via billing-service for COMPLETED appointments and marks PAID bills

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$root = Split-Path -Parent $MyInvocation.MyCommand.Definition
$dataDir = Join-Path $root '..\hms_seed_data (1)'

# Service URLs (update if you changed ports)
$userSvc = 'http://localhost:8081'
$apptSvc = 'http://localhost:8082'
$billingSvc = 'http://localhost:8083'
$doctorSvc = 'http://localhost:8085'

# Helper: convert CSV datetime to ISO-8601 UTC (assumes CSV timestamps are local-ish)
function To-ISO($s) {
    if ([string]::IsNullOrWhiteSpace($s)) { return $null }
    try {
        $dt = [DateTime]::Parse($s)
        return $dt.ToUniversalTime().ToString('o')
    } catch {
        return $s
    }
}

# Wait for service health
function Wait-Healthy($url) {
    Write-Output "Checking health for $url"
    for ($i=0; $i -lt 30; $i++) {
        try {
            $r = Invoke-RestMethod -Uri ("$url/health") -Method GET -TimeoutSec 3
            if ($r -ne $null) { Write-Output "$url healthy"; return }
        } catch {
            Start-Sleep -Seconds 1
        }
    }
    throw "Service $url not healthy"
}

Wait-Healthy $userSvc
Wait-Healthy $apptSvc
Wait-Healthy $billingSvc
Wait-Healthy $doctorSvc

# Load and create patients
$patientsCsv = Join-Path $dataDir 'hms_patients.csv'
$patients = Import-Csv -Path $patientsCsv
$patientMap = @{}
Write-Output "Seeding $($patients.Count) patients..."
foreach ($p in $patients) {
    $body = @{ name = $p.name; email = $p.email; phone = $p.phone; role = 'PATIENT' } | ConvertTo-Json
    try {
        $resp = Invoke-RestMethod -Uri ("$userSvc/v1/users") -Method POST -Headers @{ 'Content-Type'='application/json' } -Body $body
        $patientMap[$p.patient_id] = $resp.user_id
    } catch {
        Write-Warning "Failed to create patient $($p.patient_id): $_"
    }
}
Write-Output "Patients created: $($patientMap.Count)"

# If many creates failed because data already exists, try to reconcile by fetching existing users and matching by email
try {
    $allUsers = Invoke-RestMethod -Uri ("$userSvc/v1/users") -Method GET -TimeoutSec 10
    $usersList = $null
    if ($allUsers -ne $null) {
        if ($allUsers.PSObject.Properties.Name -contains 'value') { $usersList = $allUsers.value } else { $usersList = $allUsers }
        foreach ($p in $patients) {
            $match = $usersList | Where-Object { $_.email -ieq $p.email }
            if ($match) {
                if ($match -is [System.Array]) { $patientMap[$p.patient_id] = $match[0].user_id } else { $patientMap[$p.patient_id] = $match.user_id }
            }
        }
        Write-Output "Reconciled patients from user-service: $($patientMap.Count)"
    }
} catch {
    Write-Warning "Could not fetch existing users to reconcile: $_"
}

# Load and create doctors
$doctorsCsv = Join-Path $dataDir 'hms_doctors.csv'
$doctors = Import-Csv -Path $doctorsCsv
$doctorMap = @{}
Write-Output "Seeding $($doctors.Count) doctors..."
foreach ($d in $doctors) {
    $bodyObj = @{ name = $d.name; email = $d.email; phone = $d.phone; department = $d.department; experience_years = 5 }
    $body = $bodyObj | ConvertTo-Json
    try {
        $resp = Invoke-RestMethod -Uri ("$doctorSvc/v1/doctors") -Method POST -Headers @{ 'Content-Type'='application/json' } -Body $body
        $doctorMap[$d.doctor_id] = $resp.doctor_id
    } catch {
        Write-Warning "Failed to create doctor $($d.doctor_id): $_"
    }
}
Write-Output "Doctors created: $($doctorMap.Count)"

# Reconcile doctors similarly by listing existing doctors and mapping by email
try {
    $allDocs = Invoke-RestMethod -Uri ("$doctorSvc/v1/doctors") -Method GET -TimeoutSec 10
    $docsList = $null
    if ($allDocs -ne $null) {
        if ($allDocs.PSObject.Properties.Name -contains 'value') { $docsList = $allDocs.value } else { $docsList = $allDocs }
        foreach ($d in $doctors) {
            $match = $docsList | Where-Object { $_.email -ieq $d.email }
            if ($match) {
                if ($match -is [System.Array]) { $doctorMap[$d.doctor_id] = $match[0].doctor_id } else { $doctorMap[$d.doctor_id] = $match.doctor_id }
            }
        }
        Write-Output "Reconciled doctors from doctor-service: $($doctorMap.Count)"
    }
} catch {
    Write-Warning "Could not fetch existing doctors to reconcile: $_"
}

# Load appointments and create; keep mapping from original appointment_id -> created id
$apptsCsv = Join-Path $dataDir 'hms_appointments.csv'
$appts = Import-Csv -Path $apptsCsv
$apptMap = @{}
Write-Output "Seeding $($appts.Count) appointments..."
foreach ($a in $appts) {
    # Map CSV patient/doctor ids to created UUIDs
    $csvPid = $a.patient_id
    $csvDid = $a.doctor_id
    if (-not $patientMap.ContainsKey($csvPid) -or -not $doctorMap.ContainsKey($csvDid)) {
        Write-Warning "Skipping appointment $($a.appointment_id): missing patient or doctor mapping"
        continue
    }
    $payload = @{ patient_id = $patientMap[$csvPid]; doctor_id = $doctorMap[$csvDid]; start_time = (To-ISO $a.slot_start); end_time = (To-ISO $a.slot_end) }
    if ($a.department) { $payload.department = $a.department }
    $body = $payload | ConvertTo-Json
    try {
        $resp = Invoke-RestMethod -Uri ("$apptSvc/v1/appointments") -Method POST -Headers @{ 'Content-Type'='application/json' } -Body $body
        $createdId = $resp.id
        $apptMap[$a.appointment_id] = $createdId
        # update status if needed
        $status = $a.status.Trim().ToUpper()
        if ($status -eq 'COMPLETED') {
            Invoke-RestMethod -Uri ("$apptSvc/v1/appointments/$createdId/complete") -Method POST -ErrorAction Stop
        } elseif ($status -eq 'CANCELLED') {
            Invoke-RestMethod -Uri ("$apptSvc/v1/appointments/$createdId") -Method DELETE -ErrorAction Stop
        }
    } catch {
        Write-Warning "Failed appointment create for $($a.appointment_id): $_"
    }
}
Write-Output "Appointments created: $($apptMap.Count)"

# Seed prescriptions into appointments DB (prescriptions table)
$prescCsv = Join-Path $dataDir 'hms_prescriptions.csv'
$prescs = Import-Csv -Path $prescCsv
Write-Output "Seeding prescriptions (direct DB insert) - will map csv appointment ids to created ids where possible..."
# Ensure prescriptions table exists in the appointments DB (some setups place this migration under prescription-service)
$createPrescSql = @"
CREATE TABLE IF NOT EXISTS prescriptions (
    id SERIAL PRIMARY KEY,
    appointment_id VARCHAR(50) NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT now()
);
"@
try {
        docker exec -i hospital-management-system-postgres-appointments-1 psql -U hms -d appointments_db -c $createPrescSql | Out-Null
        Write-Output "Ensured prescriptions table exists in appointments_db"
} catch {
        Write-Warning "Could not ensure prescriptions table exists: $_"
}
foreach ($pr in $prescs) {
    $csvApptKey = $pr.appointment_id.ToString().Trim()
    if (-not $apptMap.ContainsKey($csvApptKey)) {
        Write-Warning "Skipping prescription $($pr.prescription_id): appointment mapping missing (csv appointment id: $csvApptKey)"
        continue
    }
    $createdApptId = $apptMap[$csvApptKey]
    $content = "medication:$($pr.medication);dosage:$($pr.dosage);days:$($pr.days)"
    $issued = To-ISO $pr.issued_at
    # Insert into postgres-appointments DB
    $sql = "INSERT INTO prescriptions (appointment_id, content, created_at) VALUES ($createdApptId, '" + ($content -replace "'", "''") + "', '" + $issued + "');"
    try {
        docker exec -i hospital-management-system-postgres-appointments-1 psql -U hms -d appointments_db -c $sql | Out-Null
    } catch {
        Write-Warning "Failed to insert prescription $($pr.prescription_id): $_"
    }
}
Write-Output "Prescriptions seeded (best-effort)"

# Seed bills via billing API for completed appointments (map appt ids)
$billsCsv = Join-Path $dataDir 'hms_bills.csv'
$bills = Import-Csv -Path $billsCsv
Write-Output "Seeding bills via billing API (will create and mark paid where CSV says PAID)..."
foreach ($b in $bills) {
    $csvAppt = $b.appointment_id
    if (-not $apptMap.ContainsKey($csvAppt)) { continue }
    $createdApptId = $apptMap[$csvAppt]
    $payload = @{ appointment_id = $createdApptId; amount = [decimal]$b.amount }
    $body = $payload | ConvertTo-Json
    try {
        $resp = Invoke-RestMethod -Uri ("$billingSvc/v1/bills") -Method POST -Headers @{ 'Content-Type'='application/json' } -Body $body
        $createdBillId = $resp.bill_id
        if ($b.status.Trim().ToUpper() -eq 'PAID') {
            Invoke-RestMethod -Uri ("$billingSvc/v1/bills/$createdBillId/pay") -Method POST
        } elseif ($b.status.Trim().ToUpper() -eq 'VOID') {
            # mark VOID directly in billing DB
            $sql = "UPDATE bills SET status='VOID' WHERE bill_id='" + $createdBillId + "';"
            docker exec -i hospital-management-system-postgres-billing-1 psql -U hms -d billing_db -c $sql | Out-Null
        }
    } catch {
        Write-Warning ("Failed to create bill for appointment " + $csvAppt + ": " + $_)
    }
}
Write-Output "Bills seeded (best-effort)"
Write-Output "Bills seeded (best-effort)"

# Export mapping files so you have CSV->service ID mapping for demo artifacts
$outDir = Join-Path $root '..\scripts_out'
if (-not (Test-Path $outDir)) { New-Item -ItemType Directory -Path $outDir | Out-Null }
# Export the mapping hashtables to JSON files. ConvertTo-Json handles hashtables directly.
$patientMap | ConvertTo-Json | Out-File -FilePath (Join-Path $outDir 'patients_map.json') -Encoding utf8
$doctorMap | ConvertTo-Json | Out-File -FilePath (Join-Path $outDir 'doctors_map.json') -Encoding utf8
$apptMap | ConvertTo-Json | Out-File -FilePath (Join-Path $outDir 'appts_map.json') -Encoding utf8

Write-Output "Seeding complete. Mapping files written to: $outDir (patients_map.json, doctors_map.json, appts_map.json)"