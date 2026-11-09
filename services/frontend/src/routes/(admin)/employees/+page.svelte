<script lang="ts">
	import { onMount } from 'svelte';
	import { 
		Card, Table, TableHead, TableHeadCell, TableBody, TableBodyRow, TableBodyCell, 
		Button, Badge, Modal, Label, Input, Select, Alert, Spinner, Fileupload
	} from 'flowbite-svelte';
	import { auth } from '$lib/stores/auth.svelte';

	let employees = $state<any[]>([]);
	let loading = $state(true);
	let errorMsg = $state('');
	let successMsg = $state('');

	// Modal Tambah Pegawai
	let showAddModal = $state(false);
	let name = $state('');
	let nik = $state('');
	let email = $state('');
	let password = $state('password123');
	let role = $state('EMPLOYEE');

	// Modal Upload Wajah Master Pegawai
	let showFaceModal = $state(false);
	let selectedEmployee = $state<any>(null);
	let faceFile = $state<FileList | null>(null);
	let uploadLoading = $state(false);

	const BACKEND_URL = 'http://localhost:8080/api/v1';

	onMount(loadEmployees);

	async function loadEmployees() {
		loading = true;
		errorMsg = '';
		try {
			const res = await fetch(`${BACKEND_URL}/master/employees`, {
				headers: {
					'Authorization': `Bearer ${auth.token}`
				}
			});
			const json = await res.json();
			if (res.ok) {
				employees = json.data || [];
			} else {
				errorMsg = json.error || 'Gagal memuat data pegawai.';
			}
		} catch (err: any) {
			errorMsg = err.message;
		} finally {
			loading = false;
		}
	}

	async function handleAddEmployee(e: Event) {
		e.preventDefault();
		errorMsg = '';
		successMsg = '';

		try {
			const res = await fetch(`${BACKEND_URL}/master/employees`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${auth.token}`
				},
				body: JSON.stringify({
					name,
					nik,
					email,
					password,
					role,
					employment_type: 'FULL_TIME',
					allow_wfh: false
				})
			});

			const json = await res.json();
			if (!res.ok) throw new Error(json.error || 'Gagal menambahkan pegawai.');

			successMsg = `Pegawai ${name} berhasil didaftarkan!`;
			showAddModal = false;
			name = '';
			nik = '';
			email = '';
			loadEmployees();
		} catch (err: any) {
			errorMsg = err.message;
		}
	}

	function openFaceUploadModal(emp: any) {
		selectedEmployee = emp;
		faceFile = null;
		showFaceModal = true;
		errorMsg = '';
	}

	async function handleUploadFace(e: Event) {
		e.preventDefault();
		if (!selectedEmployee || !faceFile || faceFile.length === 0) {
			errorMsg = 'Pilih file foto wajah terlebih dahulu.';
			return;
		}

		uploadLoading = true;
		errorMsg = '';

		try {
			const formData = new FormData();
			formData.append('photo', faceFile[0]);

			const res = await fetch(`${BACKEND_URL}/master/employees/${selectedEmployee.id}/face-biometric`, {
				method: 'POST',
				headers: {
					'Authorization': `Bearer ${auth.token}`
				},
				body: formData
			});

			const json = await res.json();
			if (!res.ok) throw new Error(json.error || 'Gagal mengunggah foto wajah.');

			successMsg = `Basis data wajah untuk pegawai ${selectedEmployee.name} berhasil didaftarkan!`;
			showFaceModal = false;
			loadEmployees();
		} catch (err: any) {
			errorMsg = err.message;
		} finally {
			uploadLoading = false;
		}
	}
</script>

<div class="space-y-6 max-w-5xl mx-auto">
	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 border-b border-gray-800 pb-4">
		<div>
			<h1 class="text-2xl font-bold text-white tracking-tight">Manajemen Employee (Pegawai)</h1>
			<p class="text-xs text-gray-400">Kelola data master pegawai dan basis data biometrik wajah master</p>
		</div>
		<Button color="green" size="sm" onclick={() => { showAddModal = true; errorMsg = ''; }}>
			➕ Tambah Pegawai Baru
		</Button>
	</div>

	{#if successMsg}
		<Alert color="green" class="text-xs" dismissable>
			✅ {successMsg}
		</Alert>
	{/if}

	{#if errorMsg}
		<Alert color="red" class="text-xs" dismissable>
			⚠️ {errorMsg}
		</Alert>
	{/if}

	<!-- Employees Table -->
	<Card class="bg-gray-900 border-gray-800 p-4 max-w-none shadow-xl">
		<Table hoverable class="text-xs text-left">
			<TableHead class="text-gray-400 border-b border-gray-800 bg-gray-800/40">
				<TableHeadCell>Nama & NIK</TableHeadCell>
				<TableHeadCell>Email</TableHeadCell>
				<TableHeadCell>Role</TableHeadCell>
				<TableHeadCell>Cabang</TableHeadCell>
				<TableHeadCell>Status Biometrik Wajah</TableHeadCell>
				<TableHeadCell>Aksi</TableHeadCell>
			</TableHead>
			<TableBody class="divide-y divide-gray-800">
				{#if loading}
					<TableBodyRow>
						<TableBodyCell colspan={6} class="text-center py-8 text-gray-500">
							<Spinner size="4" class="mr-2 inline" /> Memuat daftar pegawai...
						</TableBodyCell>
					</TableBodyRow>
				{:else if employees.length === 0}
					<TableBodyRow>
						<TableBodyCell colspan={6} class="text-center py-8 text-gray-500">
							Belum ada pegawai terdaftar.
						</TableBodyCell>
					</TableBodyRow>
				{:else}
					{#each employees as emp}
						<TableBodyRow class="hover:bg-gray-800/40">
							<TableBodyCell class="font-semibold text-white">
								{emp.name}
								<div class="text-[10px] text-gray-400 font-mono">{emp.nik}</div>
							</TableBodyCell>
							<TableBodyCell class="font-mono text-gray-300">{emp.email}</TableBodyCell>
							<TableBodyCell>
								<Badge color={emp.role === 'OWNER' ? 'purple' : emp.role === 'ADMIN_COMPANY' ? 'blue' : 'green'} class="text-[10px] font-mono">
									{emp.role}
								</Badge>
							</TableBodyCell>
							<TableBodyCell class="text-gray-300">{emp.branch_name || 'Kantor Pusat'}</TableBodyCell>
							<TableBodyCell>
								{#if emp.has_face_data}
									<Badge color="green" class="text-[10px]">✓ Wajah Terdaftar</Badge>
								{:else}
									<Badge color="yellow" class="text-[10px]">Belum Ada Foto Wajah</Badge>
								{/if}
							</TableBodyCell>
							<TableBodyCell>
								<Button size="xs" color="purple" onclick={() => openFaceUploadModal(emp)}>
									📷 Upload Wajah Master
								</Button>
							</TableBodyCell>
						</TableBodyRow>
					{/each}
				{/if}
			</TableBody>
		</Table>
	</Card>
</div>

<!-- Modal Tambah Pegawai -->
<Modal bind:open={showAddModal} size="md" autoclose={false} class="bg-gray-900 text-white">
	<h3 class="text-base font-bold text-white mb-4">Daftarkan Pegawai Baru</h3>
	<form onsubmit={handleAddEmployee} class="space-y-4 text-xs">
		<div>
			<Label for="emp_name" class="mb-1 text-xs text-gray-300">Nama Lengkap</Label>
			<Input id="emp_name" bind:value={name} required placeholder="cth: Siti Rahma" class="bg-gray-800 text-xs" />
		</div>
		<div class="grid grid-cols-2 gap-2">
			<div>
				<Label for="emp_nik" class="mb-1 text-xs text-gray-300">Nomor Induk Karyawan (NIK)</Label>
				<Input id="emp_nik" bind:value={nik} required placeholder="EMP005" class="bg-gray-800 text-xs" />
			</div>
			<div>
				<Label for="emp_role" class="mb-1 text-xs text-gray-300">Role Pengguna</Label>
				<Select id="emp_role" bind:value={role} class="bg-gray-800 text-xs">
					<option value="EMPLOYEE">Employee (Karyawan)</option>
					<option value="ADMIN_COMPANY">Admin Company (HR)</option>
				</Select>
			</div>
		</div>
		<div>
			<Label for="emp_email" class="mb-1 text-xs text-gray-300">Email Perusahaan</Label>
			<Input id="emp_email" type="email" bind:value={email} required placeholder="siti@perusahaan.com" class="bg-gray-800 text-xs" />
		</div>
		<div>
			<Label for="emp_pwd" class="mb-1 text-xs text-gray-300">Kata Sandi Awal</Label>
			<Input id="emp_pwd" type="password" bind:value={password} required class="bg-gray-800 text-xs" />
		</div>
		<div class="flex justify-end gap-2 pt-2 border-t border-gray-800">
			<Button color="light" size="sm" onclick={() => showAddModal = false}>Batal</Button>
			<Button type="submit" color="green" size="sm" class="font-bold text-slate-950">Simpan Pegawai</Button>
		</div>
	</form>
</Modal>

<!-- Modal Upload Wajah Master Pegawai (Admin) -->
<Modal bind:open={showFaceModal} size="md" autoclose={false} class="bg-gray-900 text-white">
	<h3 class="text-base font-bold text-white mb-2">Upload Basis Data Wajah Pegawai</h3>
	{#if selectedEmployee}
		<p class="text-xs text-gray-400 mb-4">
			Mendaftarkan master wajah biometrik untuk: <strong class="text-emerald-400">{selectedEmployee.name}</strong> ({selectedEmployee.nik})
		</p>
	{/if}

	<form onsubmit={handleUploadFace} class="space-y-4 text-xs">
		<div>
			<Label class="mb-2 text-xs text-gray-300">Pilih Foto Wajah Pegawai (Format JPG/PNG)</Label>
			<Fileupload bind:files={faceFile} accept="image/*" required class="bg-gray-800 text-xs" />
			<p class="text-[10px] text-gray-500 mt-1">Pastikan wajah menghadap lurus ke kamera, tidak memakai kacamata hitam/masker, dan pencahayaan terang.</p>
		</div>

		<div class="flex justify-end gap-2 pt-2 border-t border-gray-800">
			<Button color="light" size="sm" onclick={() => showFaceModal = false}>Batal</Button>
			<Button type="submit" color="purple" size="sm" disabled={uploadLoading}>
				{#if uploadLoading}
					<Spinner size="4" class="mr-2" /> Mengekstrak Biometrik...
				{:else}
					Ekstrak & Daftarkan Wajah
				{/if}
			</Button>
		</div>
	</form>
</Modal>
