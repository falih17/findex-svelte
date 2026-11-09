-- ==========================================================
-- 02: DATABASE SCHEMA & SEED DATA (MULTI-TENANT SAAS)
-- ==========================================================

-- ----------------------------------------------------------
-- 1. COMPANIES & TENANT CONFIGURATIONS
-- ----------------------------------------------------------
CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,         -- Subdomain or company code (e.g. acme-corp)
    name VARCHAR(255) NOT NULL,
    plan_tier VARCHAR(50) DEFAULT 'STANDARD', -- STARTER, STANDARD, ENTERPRISE
    max_employees INT DEFAULT 100,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS company_settings (
    company_id UUID PRIMARY KEY REFERENCES companies(id) ON DELETE CASCADE,
    timezone VARCHAR(50) DEFAULT 'Asia/Jakarta',
    face_similarity_threshold NUMERIC(3,2) DEFAULT 0.80, -- Range 0.00 to 1.00
    allow_wfh_default BOOLEAN DEFAULT FALSE,
    strict_geofence BOOLEAN DEFAULT TRUE,
    auto_overtime_approval BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ----------------------------------------------------------
-- 2. BRANCHES & GEOFENCING LOCATIONS (POSTGIS)
-- ----------------------------------------------------------
CREATE TABLE IF NOT EXISTS branches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    address TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    radius_meters INT DEFAULT 50, -- Geofencing radius in meters
    -- Auto-generated PostGIS geography point (WGS 84 SRID 4326)
    location_geom GEOGRAPHY(POINT, 4326) GENERATED ALWAYS AS (
        ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)::geography
    ) STORED,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_branches_company ON branches(company_id);
CREATE INDEX IF NOT EXISTS idx_branches_geom ON branches USING GIST(location_geom);

-- ----------------------------------------------------------
-- 3. DEPARTMENTS & EMPLOYEES
-- ----------------------------------------------------------
CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    department_id UUID REFERENCES departments(id) ON DELETE SET NULL,
    default_branch_id UUID REFERENCES branches(id) ON DELETE SET NULL,
    nik VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'EMPLOYEE', -- SUPERADMIN, HR_ADMIN, MANAGER, EMPLOYEE
    employment_type VARCHAR(50) DEFAULT 'FULL_TIME', -- FULL_TIME, CONTRACT, INTERN
    allow_wfh BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_company_nik UNIQUE (company_id, nik),
    CONSTRAINT uq_company_email UNIQUE (company_id, email)
);

CREATE INDEX IF NOT EXISTS idx_employees_lookup ON employees(company_id, is_active);

-- ----------------------------------------------------------
-- 4. BIOMETRIC MASTER DATA (PGVECTOR)
-- ----------------------------------------------------------
CREATE TABLE IF NOT EXISTS employee_biometrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL UNIQUE REFERENCES employees(id) ON DELETE CASCADE,
    master_image_url TEXT NOT NULL,
    -- 512-dimensional vector embedding extracted from ArcFace / InsightFace
    face_embedding vector(512) NOT NULL,
    liveness_score_sample NUMERIC(4,3),
    registered_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_biometrics_company_emp ON employee_biometrics(company_id, employee_id);

-- ----------------------------------------------------------
-- 5. SHIFTS & SCHEDULES
-- ----------------------------------------------------------
CREATE TABLE IF NOT EXISTS shifts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,               -- Cth: Shift Reguler Pagi, Shift Siang
    start_time TIME NOT NULL,                  -- 08:00:00 (Jam Mulai Kerja)
    check_in_deadline TIME DEFAULT '08:30:00', -- Batas Akhir Check-In (setelah ini = Terlambat/Late)
    earliest_checkout_time TIME DEFAULT '16:30:00', -- Waktu Awal Boleh Check-Out (sebelum ini = Pulang Cepat)
    end_time TIME NOT NULL,                    -- 17:00:00 (Jam Selesai Kerja)
    grace_period_late_minutes INT DEFAULT 15,
    min_work_hours NUMERIC(4,2) DEFAULT 8.0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS employee_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    shift_id UUID NOT NULL REFERENCES shifts(id),
    schedule_date DATE NOT NULL,
    is_off_day BOOLEAN DEFAULT FALSE,
    assigned_branch_id UUID REFERENCES branches(id),
    CONSTRAINT uq_emp_schedule UNIQUE (company_id, employee_id, schedule_date)
);

CREATE INDEX IF NOT EXISTS idx_schedules_emp_date ON employee_schedules(employee_id, schedule_date);

-- ----------------------------------------------------------
-- 6. ATTENDANCES (PARTITIONED BY RANGE OF DATE)
-- ----------------------------------------------------------
CREATE TABLE IF NOT EXISTS attendances (
    id UUID DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id),
    attendance_date DATE NOT NULL,
    
    -- Check-in
    check_in_time TIMESTAMPTZ,
    check_in_photo_url TEXT,
    check_in_lat DOUBLE PRECISION,
    check_in_lng DOUBLE PRECISION,
    check_in_location_geom GEOGRAPHY(POINT, 4326),
    check_in_distance_meters NUMERIC(8,2),
    check_in_similarity_score NUMERIC(4,3),
    check_in_status VARCHAR(30), -- ON_TIME, LATE, WFH, OUT_OF_GEOFENCE
    
    -- Check-out
    check_out_time TIMESTAMPTZ,
    check_out_photo_url TEXT,
    check_out_lat DOUBLE PRECISION,
    check_out_lng DOUBLE PRECISION,
    check_out_location_geom GEOGRAPHY(POINT, 4326),
    check_out_distance_meters NUMERIC(8,2),
    check_out_similarity_score NUMERIC(4,3),
    check_out_status VARCHAR(30), -- NORMAL, EARLY_LEAVE
    
    -- Summary
    work_duration_minutes INT,
    is_anomaly BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (attendance_date, id)
) PARTITION BY RANGE (attendance_date);

-- Pre-create partitions for 2026 & 2027
CREATE TABLE IF NOT EXISTS attendances_y2026m11 PARTITION OF attendances
    FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');
CREATE TABLE IF NOT EXISTS attendances_y2026m12 PARTITION OF attendances
    FOR VALUES FROM ('2026-12-01') TO ('2027-01-01');
CREATE TABLE IF NOT EXISTS attendances_y2027m01 PARTITION OF attendances
    FOR VALUES FROM ('2027-01-01') TO ('2027-02-01');

CREATE INDEX IF NOT EXISTS idx_attendances_comp_emp_date 
    ON attendances(company_id, employee_id, attendance_date);

-- ----------------------------------------------------------
-- 7. HR MODULES: LEAVE, OVERTIME, CLAIMS
-- ----------------------------------------------------------
CREATE TABLE IF NOT EXISTS leave_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL, -- Cuti Tahunan, Sakit, Melahirkan
    default_quota_days INT DEFAULT 12,
    is_deduct_annual BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS leave_balances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    leave_type_id UUID NOT NULL REFERENCES leave_types(id) ON DELETE CASCADE,
    year INT NOT NULL,
    total_quota INT NOT NULL,
    used_quota INT DEFAULT 0,
    CONSTRAINT uq_leave_balance UNIQUE (company_id, employee_id, leave_type_id, year)
);

CREATE TABLE IF NOT EXISTS leave_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    leave_type_id UUID NOT NULL REFERENCES leave_types(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_days INT NOT NULL,
    reason TEXT NOT NULL,
    attachment_url TEXT,
    status VARCHAR(30) DEFAULT 'PENDING', -- PENDING, APPROVED, REJECTED, CANCELLED
    approved_by UUID REFERENCES employees(id),
    approval_notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS overtime_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    overtime_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    estimated_hours NUMERIC(4,2) NOT NULL,
    task_description TEXT NOT NULL,
    status VARCHAR(30) DEFAULT 'PENDING',
    approved_by UUID REFERENCES employees(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reimbursement_claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    category VARCHAR(100) NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    claim_date DATE NOT NULL,
    description TEXT,
    receipt_url TEXT NOT NULL,
    status VARCHAR(30) DEFAULT 'PENDING',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ----------------------------------------------------------
-- 8. INITIAL SEED DATA FOR QUICK DEVELOPMENT
-- ----------------------------------------------------------
DO $$
DECLARE
    v_company_id UUID := '11111111-1111-1111-1111-111111111111';
    v_branch_id UUID := '22222222-2222-2222-2222-222222222222';
    v_dept_id UUID := '33333333-3333-3333-3333-333333333333';
    v_shift_id UUID := '44444444-4444-4444-4444-444444444444';
    v_admin_emp_id UUID := '55555555-5555-5555-5555-555555555555';
BEGIN
    -- 1. Demo Company
    INSERT INTO companies (id, code, name, plan_tier, max_employees)
    VALUES (v_company_id, 'demo-corp', 'PT Findex Demo Indonesia', 'ENTERPRISE', 500)
    ON CONFLICT (id) DO NOTHING;

    INSERT INTO company_settings (company_id, timezone, face_similarity_threshold, allow_wfh_default, strict_geofence)
    VALUES (v_company_id, 'Asia/Jakarta', 0.80, TRUE, TRUE)
    ON CONFLICT (company_id) DO NOTHING;

    -- 2. Demo Branch (Kantor Pusat Sudirman Jakarta: -6.2088, 106.8456)
    INSERT INTO branches (id, company_id, name, address, latitude, longitude, radius_meters)
    VALUES (v_branch_id, v_company_id, 'Headquarters Jakarta', 'Jl. Jend. Sudirman Kav. 52-53, Jakarta Selatan', -6.2088, 106.8456, 100)
    ON CONFLICT (id) DO NOTHING;

    -- 3. Demo Department
    INSERT INTO departments (id, company_id, name)
    VALUES (v_dept_id, v_company_id, 'Engineering & IT')
    ON CONFLICT (id) DO NOTHING;

    -- 4. Demo Shift (08:00 - 17:00 WIB, Batas Check-in: 08:30, Awal Checkout: 16:30)
    INSERT INTO shifts (id, company_id, name, start_time, check_in_deadline, earliest_checkout_time, end_time, grace_period_late_minutes)
    VALUES (v_shift_id, v_company_id, 'Shift Reguler Pagi', '08:00:00', '08:30:00', '16:30:00', '17:00:00', 15)
    ON CONFLICT (id) DO NOTHING;

    -- 5. Demo Users with 3 Explicit Roles (Owner, Admin Company, Employee)
    -- A. Platform Owner (Bisa akses semua modul)
    INSERT INTO employees (id, company_id, nik, name, email, password_hash, role, allow_wfh)
    VALUES ('00000000-0000-0000-0000-000000000001', v_company_id, 'OWNER01', 'SaaS Platform Owner', 'owner@findex.local', crypt('password123', gen_salt('bf', 10)), 'OWNER', TRUE)
    ON CONFLICT (company_id, email) DO NOTHING;

    -- B. Admin Company (Akses: Employee, Jadwal, Lokasi Absen, Laporan)
    INSERT INTO employees (id, company_id, department_id, default_branch_id, nik, name, email, password_hash, role, allow_wfh)
    VALUES (v_admin_emp_id, v_company_id, v_dept_id, v_branch_id, 'EMP001', 'Admin Findex', 'admin@findex.local', crypt('password123', gen_salt('bf', 10)), 'ADMIN_COMPANY', TRUE)
    ON CONFLICT (company_id, email) DO NOTHING;

    -- C. Employee (Presensi Wajah & GPS, Riwayat, Upload Wajah Mandiri)
    INSERT INTO employees (id, company_id, department_id, default_branch_id, nik, name, email, password_hash, role, allow_wfh)
    VALUES ('66666666-6666-6666-6666-666666666666', v_company_id, v_dept_id, v_branch_id, 'EMP002', 'Budi Pratama', 'employee@findex.local', crypt('password123', gen_salt('bf', 10)), 'EMPLOYEE', FALSE)
    ON CONFLICT (company_id, email) DO NOTHING;

    -- 6. Demo Leave Types
    INSERT INTO leave_types (company_id, name, default_quota_days, is_deduct_annual)
    VALUES 
        (v_company_id, 'Cuti Tahunan', 12, TRUE),
        (v_company_id, 'Cuti Sakit', 14, FALSE),
        (v_company_id, 'Cuti Melahirkan', 90, FALSE)
    ON CONFLICT DO NOTHING;

    RAISE NOTICE 'Demo data seeded successfully for company: PT Findex Demo Indonesia (ID: %)', v_company_id;
END $$;
