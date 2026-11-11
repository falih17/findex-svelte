<script lang="ts">
	import { onMount } from 'svelte';
	import { 
		Modal, Label, Input, Select, Alert, Spinner, Fileupload
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
					role
				})
			});

			const json = await res.json();
			if (!res.ok) {
				throw new Error(json.error || 'Gagal menambahkan pegawai.');
			}

			successMsg = `Pegawai "${name}" berhasil didaftarkan!`;
			showAddModal = false;
			name = '';
			nik = '';
			email = '';
			loadEmployees();
		} catch (err: any) {
			errorMsg = err.message;
		}
	}

	function openFaceUpload(emp: any) {
		selectedEmployee = emp;
		faceFile = null;
		showFaceModal = true;
	}

	async function handleUploadFace(e: Event) {
		e.preventDefault();
		if (!faceFile || faceFile.length === 0 || !selectedEmployee) {
			errorMsg = 'Pilih file foto wajah terlebih dahulu!';
			return;
		}

		uploadLoading = true;
		errorMsg = '';
		successMsg = '';

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
			if (!res.ok) {
				throw new Error(json.error || 'Gagal mengekstraksi biometrik wajah.');
			}

			successMsg = `Wajah master biometrik untuk ${selectedEmployee.name} berhasil disimpan ke AI Face Service!`;
			showFaceModal = false;
			loadEmployees();
		} catch (err: any) {
			errorMsg = err.message;
		} finally {
			uploadLoading = false;
		}
	}
</script>

<div class="space-y-6">
	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pb-4 border-b border-slate-200">
		<div>
			<h1 class="text-2xl font-bold text-slate-900 tracking-tight">Karyawan & Divisi</h1>
			<p class="text-xs text-slate-500 mt-0.5">Kelola data pegawai, hak akses, dan basis data biometrik wajah</p>
		</div>
		<button
			type="button"
			onclick={() => { showAddModal = true; errorMsg = ''; successMsg = ''; }}
			class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold text-xs rounded-lg shadow-sm shadow-indigo-600/20 transition-all"
		>
			<span>➕ Tambah Pegawai Baru</span>
		</button>
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

	<!-- Employees Table Card -->
	<div class="bg-white border border-slate-200/90 rounded-2xl p-6 shadow-2xs space-y-4">
		{#if loading}
			<div class="text-center py-12 text-slate-500 text-sm">
				<Spinner size="4" class="mr-2 inline" /> Memuat data pegawai...
			</div>
		{:else if employees.length === 0}
			<div class="text-center py-10 text-slate-400 text-xs">
				Belum ada data pegawai. Silakan tambahkan pegawai pertama.
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full text-left text-xs">
					<thead>
						<tr class="text-[11px] font-bold text-slate-400 uppercase tracking-wider border-b border-slate-100">
							<th class="pb-3 pl-2">Nama Pegawai</th>
							<th class="pb-3">NIK</th>
							<th class="pb-3">Email Akun</th>
							<th class="pb-3">Role Akses</th>
							<th class="pb-3">Status Biometrik</th>
							<th class="pb-3 pr-2 text-right">Kelola Wajah</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100">
						{#each employees as emp}
							<tr class="hover:bg-slate-50/70 transition-colors">
								<td class="py-3.5 pl-2 font-bold text-slate-900">{emp.name}</td>
								<td class="py-3.5 font-mono text-slate-500">{emp.nik || '-'}</td>
								<td class="py-3.5 text-slate-600">{emp.email}</td>
								<td class="py-3.5">
									<span class={`px-2 py-0.5 rounded-full font-bold text-[10px] ${
										emp.role === 'OWNER'
											? 'bg-purple-50 text-purple-700 border border-purple-200'
											: emp.role === 'ADMIN_COMPANY'
											? 'bg-indigo-50 text-indigo-700 border border-indigo-200'
											: 'bg-slate-100 text-slate-700'
									}`}>
										{emp.role}
									</span>
								</td>
								<td class="py-3.5">
									{#if emp.face_embedding}
										<span class="px-2 py-0.5 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-700 font-bold text-[10px] flex items-center gap-1 w-max">
											<span>✓</span> Terdaftar (pgvector)
										</span>
									{:else}
										<span class="px-2 py-0.5 rounded-full bg-amber-50 border border-amber-200 text-amber-700 font-bold text-[10px] w-max">
											Belum Terdaftar
										</span>
									{/if}
								</td>
								<td class="py-3.5 pr-2 text-right">
									<button
										type="button"
										onclick={() => openFaceUpload(emp)}
										class="px-2.5 py-1 bg-white hover:bg-slate-50 border border-slate-200 text-indigo-600 hover:text-indigo-700 text-xs font-semibold rounded-lg shadow-2xs transition-colors"
									>
										📸 Upload Wajah
									</button>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>

<!-- Modal Tambah Pegawai -->
<Modal bind:open={showAddModal} size="md" autoclose={false} title="Daftarkan Pegawai Baru">
	<form onsubmit={handleAddEmployee} class="space-y-4 text-xs">
		<div>
			<Label for="ename" class="mb-1 text-slate-700">Nama Lengkap</Label>
			<Input id="ename" bind:value={name} placeholder="Contoh: Budi Santoso" required />
		</div>

		<div>
			<Label for="enik" class="mb-1 text-slate-700">Nomor Induk Karyawan (NIK)</Label>
			<Input id="enik" bind:value={nik} placeholder="Contoh: 2024001" required />
		</div>

		<div>
			<Label for="eemail" class="mb-1 text-slate-700">Email Akun</Label>
			<Input id="eemail" type="email" bind:value={email} placeholder="budi@perusahaan.com" required />
		</div>

		<div>
			<Label for="epass" class="mb-1 text-slate-700">Kata Sandi Awal</Label>
			<Input id="epass" type="text" bind:value={password} required />
		</div>

		<div>
			<Label for="erole" class="mb-1 text-slate-700">Role Pengguna</Label>
			<Select id="erole" bind:value={role}>
				<option value="EMPLOYEE">EMPLOYEE (Hanya Presensi & Riwayat)</option>
				<option value="ADMIN_COMPANY">ADMIN_COMPANY (Kelola Employee, Jadwal, Lokasi)</option>
			</Select>
		</div>

		<div class="flex justify-end gap-2 pt-2 border-t border-slate-100">
			<button type="button" class="px-3 py-1.5 rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50 text-xs" onclick={() => showAddModal = false}>
				Batal
			</button>
			<button type="submit" class="px-4 py-1.5 rounded-lg bg-indigo-600 hover:bg-indigo-700 text-white font-semibold text-xs shadow-2xs">
				Simpan Pegawai
			</button>
		</div>
	</form>
</Modal>

<!-- Modal Upload Wajah Master Pegawai (Oleh Admin) -->
<Modal bind:open={showFaceModal} size="md" autoclose={false} title={`Upload Wajah Master Biometrik - ${selectedEmployee?.name}`}>
	<form onsubmit={handleUploadFace} class="space-y-4 text-xs">
		<p class="text-slate-500 text-xs">
			Unggah foto wajah resmi beresolusi jelas menghadap ke depan tanpa penutup wajah (kacamata hitam / masker). AI Service akan mengekstraksi 512-dimensi vektor dan menyimpannya ke pgvector.
		</p>

		<div>
			<Label for="facefile" class="mb-1 text-slate-700">Pilih Foto JPEG / PNG</Label>
			<Fileupload id="facefile" bind:files={faceFile} accept="image/jpeg,image/png" required />
		</div>

		<div class="flex justify-end gap-2 pt-2 border-t border-slate-100">
			<button type="button" class="px-3 py-1.5 rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50 text-xs" onclick={() => showFaceModal = false}>
				Batal
			</button>
			<button type="submit" disabled={uploadLoading} class="px-4 py-1.5 rounded-lg bg-indigo-600 hover:bg-indigo-700 text-white font-semibold text-xs shadow-2xs">
				{uploadLoading ? 'Mengekstrak AI Vektor...' : 'Ekstrak & Simpan Vektor Wajah'}
			</button>
		</div>
	</form>
</Modal>