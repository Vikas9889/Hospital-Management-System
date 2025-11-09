# Demo Flow Script - Creates a complete appointment cycle
# Usage: ./scripts/demo-flow.ps1

Write-Output "Creating demo resources..."

# 1. Create a new patient (through user service)
Write-Output "`n1. Creating new patient:"
$user = Invoke-RestMethod -Uri "http://localhost:8081/v1/users" -Method Post -Body (@{
    name = "Demo Patient Bits"
    email = "demo.patientbits@example.com"
    phone = "9999988886"
    role = "PATIENT"
} | ConvertTo-Json) -ContentType "application/json"
$user | ConvertTo-Json

# 2. Create a new doctor (through doctor service)
Write-Output "`n2. Creating new doctor:"
$doctor = Invoke-RestMethod -Uri "http://localhost:8085/v1/doctors" -Method Post -Body (@{
    name = "Dr. Demo Doctor Bits"
    email = "demo.doctorbits@example.com"
    specialization = "General Medicine"
} | ConvertTo-Json) -ContentType "application/json"
Write-Output "`nCreated Doctor:"
$doctor | ConvertTo-Json

# 3. Create an appointment
Write-Output "`n3. Creating new appointment:"
$startTime = (Get-Date).AddDays(1)
$endTime = $startTime.AddHours(1)

$appointment = Invoke-RestMethod -Uri "http://localhost:8082/v1/appointments" -Method Post -Body (@{
    patient_id = $user.user_id
    doctor_id = $doctor.doctor_id
    start_time = $startTime.ToString('o')
    end_time = $endTime.ToString('o')
} | ConvertTo-Json) -ContentType "application/json"
$appointment | ConvertTo-Json

# 4. Mark appointment as completed
Write-Output "`n4. Completing appointment:"
$updateAppointment = Invoke-RestMethod -Uri "http://localhost:8082/v1/appointments/$($appointment.id)/complete" -Method Post -ContentType "application/json"
$updateAppointment | ConvertTo-Json

# 5. Create a prescription
Write-Output "`n5. Creating prescription:"
$prescription = Invoke-RestMethod -Uri "http://localhost:8087/v1/prescriptions" -Method Post -Body (@{
    appointment_id = $appointment.id
    doctor_id = $doctor.doctor_id
    patient_id = $user.user_id
    medication = "Demo Medicine 2025"
    dosage = "2 tablets"
    frequency = "thrice daily"
    duration = "3 days"
    notes = "Take with warm water"
} | ConvertTo-Json) -ContentType "application/json"
$prescription | ConvertTo-Json

# 6. Generate bill
Write-Output "`n6. Generating bill:"
$bill = Invoke-RestMethod -Uri "http://localhost:8083/v1/bills" -Method Post -Body (@{
    appointment_id = $appointment.id
    amount = 1000
    currency = "INR"
    status = "PENDING"
} | ConvertTo-Json) -ContentType "application/json"
$bill | ConvertTo-Json

# 7. Make payment
Write-Output "`n7. Making payment:"
$payment = Invoke-RestMethod -Uri "http://localhost:8088/v1/payments" -Method Post -Body (@{
    bill_id = $bill.bill_id
    amount = $bill.amount
    method = "CARD"
    status = "COMPLETED"
} | ConvertTo-Json) -ContentType "application/json"
$payment | ConvertTo-Json

# 8. Verify the flow by retrieving all created resources
Write-Output "`n8. Verifying all resources:"
Write-Output "`n8.1. Patient Details:"
Invoke-RestMethod -Uri "http://localhost:8081/v1/users/$($user.user_id)" -Method Get | ConvertTo-Json

Write-Output "`n8.2. Doctor Details:"
Invoke-RestMethod -Uri "http://localhost:8085/v1/doctors/$($doctor.doctor_id)" -Method Get | ConvertTo-Json

Write-Output "`n8.3. Appointment Details:"
Invoke-RestMethod -Uri "http://localhost:8082/v1/appointments/$($appointment.id)" -Method Get | ConvertTo-Json

Write-Output "`n8.4. Prescription Details:"
Invoke-RestMethod -Uri "http://localhost:8087/v1/prescriptions?appointment_id=$($appointment.id)" -Method Get | ConvertTo-Json

Write-Output "`n8.5. Bill Details:"
Invoke-RestMethod -Uri "http://localhost:8083/v1/bills/$($bill.bill_id)" -Method Get | ConvertTo-Json

Write-Output "`n8.6. Payment Details:"
Invoke-RestMethod -Uri "http://localhost:8088/v1/payments?bill_id=$($bill.bill_id)" -Method Get | ConvertTo-Json