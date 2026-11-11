<script lang="ts">
	import { onMount } from 'svelte';
	import { Modal, Label, Input, Select, Alert, Spinner } from 'flowbite-svelte';
	import { auth } from '$lib/stores/auth.svelte';

	let companies = $state<any[]>([]);
	let loading = $state(true);
	let errorMsg = $state('');
	let successMsg = $state('');

	// Filters
	let searchQuery = $state('');
	let statusFilter = $state<'ALL' | 'ACTIVE' | 'INACTIVE'>('ALL');
	let tierFilter = $state('ALL');

	// Modal Tambah Tenant
	let showAddModal = $state(false);
	let newCode = $state('');
	let newName = $state('');
	let newTier = $state('STANDARD');
	let newMaxEmp = $state(100);
	let newIsActive = $state(true);

	// Modal Edit Tenant
	let showEditModal = $state(false);
	let editCompany = $state<any>(null);

	// Modal Konfirmasi Toggle Status
	let showToggleModal = $state(false);
	let targetCompany = $state<any>(null);
	let toggleLoading = $state(false);

	const BACKEND_URL = 'http://localhost:8080/api/v1';

	onMount(loadCompanies);

	async function loadCompanies() {
		loading = true;
		errorMsg = '';
		try {
			if (!auth.token) {
				errorMsg = 'Sesi Anda telah berakhir. Silakan login kembali sebagai OWNER.';
				return;
			}

			const res = await fetch(`${BACKEND_URL}/master/companies`, {
				headers: {
					Authorization: `Bearer ${auth.token}`
				}
			});

			const json = await res.json();
			if (res.ok) {
				companies = json.data || [];
			} else {
				errorMsg = json.error || 'Gagal memuat daftar tenant/company.';
			}
		} catch (err: any) {
			errorMsg = err.message;
		} finally {
			loading = false;
		}
	}

	// Filtered Companies Computed
	let filteredCompanies = $derived(
		companies.filter((c) => {
			const matchSearch =
				c.name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
				c.code?.toLowerCase().includes(searchQuery.toLowerCase());

			const matchStatus =
				statusFilter === 'ALL'
					? true
					: statusFilter === 'ACTIVE'
					? c.is_active === true
					: c.is_active === false;

			const matchTier = tierFilter === 'ALL' ? true : c.plan_tier === tierFilter;

			return matchSearch && matchStatus && matchTier;
		})
	);

	// Summary KPI stats
	let stats = $derived({
		total: companies.length,
		active: companies.filter((c) => c.is_active).length,
		inactive: companies.filter((c) => !c.is_active).length,
		totalEmployees: companies.reduce((acc, c) => acc + (c.total_employees || 0), 0)
	});

	async function handleCreateCompany(e: Event) {
		e.preventDefault();
		errorMsg = '';
		successMsg = '';

		try {
			const res = await fetch(`${BACKEND_URL}/master/companies`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${auth.token}`
				},
				body: JSON.stringify({
					code: newCode.trim().toLowerCase(),
					name: newName.trim(),
					plan_tier: newTier,
					max_employees: Number(newMaxEmp),
					is_active: newIsActive
				})
			});

			const json = await res.json();
			if (!res.ok) {
				throw new Error(json.error || 'Gagal mendaftarkan tenant.');
			}

			successMsg = `Tenant "${newName}" berhasil didaftarkan!`;
			showAddModal = false;
			newCode = '';
			newName = '';
			newTier = 'STANDARD';
			newMaxEmp = 100;
			newIsActive = true;
			loadCompanies();
		} catch (err: any) {
			errorMsg = err.message;
		}
	}

	function openEdit(c: any) {
		editCompany = { ...c };
		showEditModal = true;
	}

	async function handleUpdateCompany(e: Event) {
		e.preventDefault();
		if (!editCompany) return;
		errorMsg = '';
		successMsg = '';

		try {
			const res = await fetch(`${BACKEND_URL}/master/companies/${editCompany.id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${auth.token}`
				},
				body: JSON.stringify({
					name: editCompany.name,
					plan_tier: editCompany.plan_tier,
					max_employees: Number(editCompany.max_employees),
					is_active: editCompany.is_active
				})
			});

			const json = await res.json();
			if (!res.ok) {
				throw new Error(json.error || 'Gagal memperbarui tenant.');
			}

			successMsg = `Data tenant "${editCompany.name}" berhasil diperbarui!`;
			showEditModal = false;
			loadCompanies();
		} catch (err: any) {
			errorMsg = err.message;
		}
	}

	function openToggleStatus(c: any) {
		targetCompany = c;
		showToggleModal = true;
	}

	async function confirmToggleStatus() {
		if (!targetCompany) return;
		toggleLoading = true;
		errorMsg = '';
		successMsg = '';

		try {
			const nextStatus = !targetCompany.is_active;
			const res = await fetch(`${BACKEND_URL}/master/companies/${targetCompany.id}/status`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${auth.token}`
				},
				body: JSON.stringify({ is_active: nextStatus })
			});

			const json = await res.json();
			if (!res.ok) {
				throw new Error(json.error || 'Gagal mengubah status tenant.');
			}

			successMsg = json.message || `Status tenant "${targetCompany.name}" berhasil diubah!`;
			showToggleModal = false;
			loadCompanies();
		} catch (err: any) {
			errorMsg = err.message;
		} finally {
			toggleLoading = false;
		}
	}

	function autoFillSlug() {
		if (newName && !newCode) {
			newCode = newName
				.toLowerCase()
				.replace(/[^a-z0-9]+/g, '-')
				.replace(/(^-|-$)/g, '');
		}
	}
</script>

<div class="space-y-6">
	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pb-4 border-b border-slate-200">
		<div>
			<div class="flex items-center gap-2">
				<h1 class="text-2xl font-bold text-slate-900 tracking-tight">Manajemen Tenant & Company</h1>
				<span class="px-2.5 py-0.5 rounded-full bg-purple-50 border border-purple-200 text-purple-700 font-bold text-[11px] font-mono">
					👑 OWNER ACCESS
				</span>
			</div>
			<p class="text-xs text-slate-500 mt-0.5">
				Kelola pendaftaran multi-tenant, kuota pegawai, serta aktivasi dan penonaktifan akses perusahaan.
			</p>
		</div>

		{#if auth.user?.role === 'OWNER'}
			<button
				type="button"
				onclick={() => { showAddModal = true; errorMsg = ''; successMsg = ''; }}
				class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold text-xs rounded-lg shadow-sm shadow-indigo-600/20 transition-all"
			>
				<span>➕ Daftarkan Tenant Baru</span>
			</button>
		{/if}
	</div>

	{#if auth.user?.role !== 'OWNER'}
		<div class="p-4 bg-amber-50 border border-amber-200 rounded-xl text-amber-800 text-xs space-y-2">
			<div class="font-bold flex items-center gap-1.5 text-sm">
				<span>🔒</span>
				<span>Akses Terbatas Khusus Role OWNER</span>
			</div>
			<p>
				Saat ini Anda masuk sebagai <strong class="font-mono text-slate-900">{auth.user?.role || 'GUEST'}</strong>. Modul ini dikhususkan bagi pemilik platform SaaS (Owner) untuk mengontrol izin operasional seluruh tenant.
			</p>
			<div>
				<a href="/login" class="inline-block px-3 py-1.5 bg-amber-600 hover:bg-amber-700 text-white font-semibold rounded-lg text-xs transition-colors">
					Beralih ke Akun Owner (owner@findex.local)
				</a>
			</div>
		</div>
	{/if}

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

	<!-- 4 KPI Summary Cards -->
	<div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
		<div class="bg-white border border-slate-200/90 rounded-2xl p-4 shadow-2xs space-y-1">
			<div class="text-[11px] font-bold text-slate-400 uppercase tracking-wider">Total Tenant</div>
			<div class="text-3xl font-black text-slate-900">{stats.total}</div>
			<div class="text-[11px] text-slate-500">Perusahaan terdaftar di sistem</div>
		</div>

		<div class="bg-white border border-slate-200/90 rounded-2xl p-4 shadow-2xs space-y-1">
			<div class="text-[11px] font-bold text-emerald-600 uppercase tracking-wider flex items-center gap-1">
				<span class="w-2 h-2 rounded-full bg-emerald-500"></span>
				<span>Tenant Aktif</span>
			</div>
			<div class="text-3xl font-black text-emerald-600">{stats.active}</div>
			<div class="text-[11px] text-emerald-700/80">Operasional presensi berjalan</div>
		</div>

		<div class="bg-white border border-slate-200/90 rounded-2xl p-4 shadow-2xs space-y-1">
			<div class="text-[11px] font-bold text-rose-500 uppercase tracking-wider flex items-center gap-1">
				<span class="w-2 h-2 rounded-full bg-rose-500"></span>
				<span>Tenant Non-Aktif</span>
			</div>
			<div class="text-3xl font-black text-rose-600">{stats.inactive}</div>
			<div class="text-[11px] text-rose-700/80">Akses diblokir / dibekukan</div>
		</div>

		<div class="bg-white border border-slate-200/90 rounded-2xl p-4 shadow-2xs space-y-1">
			<div class="text-[11px] font-bold text-indigo-600 uppercase tracking-wider">Total Karyawan</div>
			<div class="text-3xl font-black text-indigo-600">{stats.totalEmployees}</div>
			<div class="text-[11px] text-slate-500">Akun terakumulasi seluruh tenant</div>
		</div>
	</div>

	<!-- Filter & Search Controls -->
	<div class="bg-white border border-slate-200/90 rounded-2xl p-4 shadow-2xs">
		<div class="grid grid-cols-1 sm:grid-cols-3 gap-3 text-xs">
			<!-- Search input -->
			<div>
				<label for="search" class="block mb-1 font-semibold text-slate-700">Cari Perusahaan</label>
				<input
					id="search"
					type="text"
					bind:value={searchQuery}
					placeholder="Ketik nama atau kode perusahaan..."
					class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs focus:ring-2 focus:ring-indigo-500"
				/>
			</div>

			<!-- Status filter -->
			<div>
				<label for="status" class="block mb-1 font-semibold text-slate-700">Status Aktivasi</label>
				<select
					id="status"
					bind:value={statusFilter}
					class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs"
				>
					<option value="ALL">Semua Status (Aktif & Non-Aktif)</option>
					<option value="ACTIVE">Hanya Tenant Aktif</option>
					<option value="INACTIVE">Hanya Tenant Non-Aktif</option>
				</select>
			</div>

			<!-- Tier filter -->
			<div>
				<label for="tier" class="block mb-1 font-semibold text-slate-700">Paket Langganan</label>
				<select
					id="tier"
					bind:value={tierFilter}
					class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs"
				>
					<option value="ALL">Semua Paket</option>
					<option value="STARTER">STARTER</option>
					<option value="STANDARD">STANDARD</option>
					<option value="ENTERPRISE">ENTERPRISE</option>
				</select>
			</div>
		</div>
	</div>

	<!-- Tenant Table Card -->
	<div class="bg-white border border-slate-200/90 rounded-2xl p-6 shadow-2xs space-y-4">
		{#if loading}
			<div class="text-center py-12 text-slate-500 text-sm">
				<Spinner size="4" class="mr-2 inline" /> Memuat data perusahaan...
			</div>
		{:else if filteredCompanies.length === 0}
			<div class="text-center py-12 text-slate-400 text-xs">
				Tidak ada tenant yang cocok dengan filter pencarian.
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full text-left text-xs">
					<thead>
						<tr class="text-[11px] font-bold text-slate-400 uppercase tracking-wider border-b border-slate-100">
							<th class="pb-3 pl-2">Nama Perusahaan & Kode</th>
							<th class="pb-3">Paket Langganan</th>
							<th class="pb-3">Kapasitas Karyawan</th>
							<th class="pb-3">Cabang Geofence</th>
							<th class="pb-3">Status Akses Tenant</th>
							<th class="pb-3 pr-2 text-right">Aksi Manajemen</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100">
						{#each filteredCompanies as comp}
							<tr class="hover:bg-slate-50/70 transition-colors">
								<!-- Nama Perusahaan & Kode -->
								<td class="py-3.5 pl-2">
									<div class="flex items-center gap-3">
										<div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-indigo-500 to-purple-600 text-white font-bold text-xs flex items-center justify-center shrink-0 shadow-2xs">
											{comp.name.slice(0, 2).toUpperCase()}
										</div>
										<div>
											<div class="font-bold text-slate-900 text-sm">{comp.name}</div>
											<div class="text-[11px] text-slate-400 font-mono">
												Kode: <span class="text-indigo-600 font-semibold">{comp.code}</span>
											</div>
										</div>
									</div>
								</td>

								<!-- Paket Langganan -->
								<td class="py-3.5">
									<span class={`px-2.5 py-1 rounded-full font-bold text-[10px] uppercase font-mono ${
										comp.plan_tier === 'ENTERPRISE'
											? 'bg-purple-50 text-purple-700 border border-purple-200'
											: comp.plan_tier === 'STANDARD'
											? 'bg-indigo-50 text-indigo-700 border border-indigo-200'
											: 'bg-blue-50 text-blue-700 border border-blue-200'
									}`}>
										{comp.plan_tier}
									</span>
								</td>

								<!-- Kapasitas Karyawan -->
								<td class="py-3.5">
									<div class="font-semibold text-slate-800">
										{comp.total_employees || 0} / {comp.max_employees} Pegawai
									</div>
									<div class="w-28 bg-slate-100 rounded-full h-1.5 mt-1.5 overflow-hidden">
										<div
											class={`h-1.5 rounded-full ${
												(comp.total_employees || 0) >= comp.max_employees
													? 'bg-rose-500'
													: 'bg-indigo-600'
											}`}
											style="width: {Math.min(100, Math.round(((comp.total_employees || 0) / (comp.max_employees || 1)) * 100))}%;"
										></div>
									</div>
								</td>

								<!-- Cabang Geofence -->
								<td class="py-3.5 font-medium text-slate-700">
									<span class="inline-flex items-center gap-1">
										<span>📍</span>
										<span>{comp.total_branches || 0} Cabang</span>
									</span>
								</td>

								<!-- Status Akses Tenant (Aktivasi / Non-Aktivasi) -->
								<td class="py-3.5">
									{#if comp.is_active}
										<div class="space-y-1">
											<span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-emerald-50 border border-emerald-300 text-emerald-700 font-bold text-[11px]">
												<span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
												<span>Aktif</span>
											</span>
											<div>
												<button
													type="button"
													onclick={() => openToggleStatus(comp)}
													class="text-[11px] text-rose-600 hover:text-rose-800 font-medium hover:underline block"
												>
													Matikan / Nonaktifkan
												</button>
											</div>
										</div>
									{:else}
										<div class="space-y-1">
											<span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-rose-50 border border-rose-300 text-rose-700 font-bold text-[11px]">
												<span class="w-2 h-2 rounded-full bg-rose-500"></span>
												<span>Non-Aktif</span>
											</span>
											<div>
												<button
													type="button"
													onclick={() => openToggleStatus(comp)}
													class="text-[11px] text-emerald-600 hover:text-emerald-800 font-medium hover:underline block"
												>
													Aktifkan Kembali
												</button>
											</div>
										</div>
									{/if}
								</td>

								<!-- Aksi Manajemen -->
								<td class="py-3.5 pr-2 text-right">
									<div class="flex items-center justify-end gap-1.5">
										<button
											type="button"
											onclick={() => openEdit(comp)}
											class="px-2.5 py-1 bg-white hover:bg-slate-50 border border-slate-200 text-slate-700 text-xs font-semibold rounded-lg shadow-2xs transition-colors"
										>
											✏️ Edit
										</button>

										<button
											type="button"
											onclick={() => openToggleStatus(comp)}
											class={`px-2.5 py-1 text-xs font-semibold rounded-lg transition-colors shadow-2xs ${
												comp.is_active
													? 'bg-rose-50 text-rose-700 border border-rose-200 hover:bg-rose-100'
													: 'bg-emerald-600 text-white hover:bg-emerald-700'
											}`}
										>
											{comp.is_active ? 'Nonaktifkan' : 'Aktifkan'}
										</button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>

<!-- Modal Tambah Tenant Baru -->
<Modal bind:open={showAddModal} size="md" autoclose={false} title="Daftarkan Tenant / Perusahaan Baru">
	<form onsubmit={handleCreateCompany} class="space-y-4 text-xs">
		<div>
			<Label for="cname" class="mb-1 text-slate-700">Nama Perusahaan / Organisasi</Label>
			<Input
				id="cname"
				bind:value={newName}
				oninput={autoFillSlug}
				placeholder="Contoh: PT Sinar Abadi Perkasa"
				required
			/>
		</div>

		<div>
			<Label for="ccode" class="mb-1 text-slate-700">Kode Tenant / Subdomain</Label>
			<Input
				id="ccode"
				bind:value={newCode}
				placeholder="sinar-abadi"
				required
			/>
			<p class="text-[10px] text-slate-400 mt-0.5">Digunakan sebagai identitas unik tenant (slug).</p>
		</div>

		<div class="grid grid-cols-2 gap-3">
			<div>
				<Label for="ctier" class="mb-1 text-slate-700">Paket Langganan</Label>
				<Select id="ctier" bind:value={newTier}>
					<option value="STARTER">STARTER</option>
					<option value="STANDARD">STANDARD</option>
					<option value="ENTERPRISE">ENTERPRISE</option>
				</Select>
			</div>

			<div>
				<Label for="cmax" class="mb-1 text-slate-700">Maksimal Pegawai</Label>
				<Input id="cmax" type="number" bind:value={newMaxEmp} min="5" required />
			</div>
		</div>

		<div>
			<Label for="cactive" class="mb-1 text-slate-700">Status Awal Tenant</Label>
			<Select id="cactive" bind:value={newIsActive}>
				<option value={true}>Langsung Aktif (Siap Digunakan)</option>
				<option value={false}>Non-Aktif (Dalam Peninjauan)</option>
			</Select>
		</div>

		<div class="flex justify-end gap-2 pt-2 border-t border-slate-100">
			<button
				type="button"
				class="px-3 py-1.5 rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50 text-xs"
				onclick={() => (showAddModal = false)}
			>
				Batal
			</button>
			<button
				type="submit"
				class="px-4 py-1.5 rounded-lg bg-indigo-600 hover:bg-indigo-700 text-white font-semibold text-xs shadow-2xs"
			>
				Simpan & Daftarkan Tenant
			</button>
		</div>
	</form>
</Modal>

<!-- Modal Edit Tenant -->
{#if editCompany}
	<Modal bind:open={showEditModal} size="md" autoclose={false} title={`Edit Tenant - ${editCompany.name}`}>
		<form onsubmit={handleUpdateCompany} class="space-y-4 text-xs">
			<div>
				<Label for="editName" class="mb-1 text-slate-700">Nama Perusahaan</Label>
				<Input id="editName" bind:value={editCompany.name} required />
			</div>

			<div class="grid grid-cols-2 gap-3">
				<div>
					<Label for="editTier" class="mb-1 text-slate-700">Paket Langganan</Label>
					<Select id="editTier" bind:value={editCompany.plan_tier}>
						<option value="STARTER">STARTER</option>
						<option value="STANDARD">STANDARD</option>
						<option value="ENTERPRISE">ENTERPRISE</option>
					</Select>
				</div>

				<div>
					<Label for="editMax" class="mb-1 text-slate-700">Maksimal Pegawai</Label>
					<Input id="editMax" type="number" bind:value={editCompany.max_employees} min="1" required />
				</div>
			</div>

			<div>
				<Label for="editActive" class="mb-1 text-slate-700">Status Akses Operasional</Label>
				<Select id="editActive" bind:value={editCompany.is_active}>
					<option value={true}>Aktif (Bisa Login & Presensi)</option>
					<option value={false}>Non-Aktif (Akses Diblokir)</option>
				</Select>
			</div>

			<div class="flex justify-end gap-2 pt-2 border-t border-slate-100">
				<button
					type="button"
					class="px-3 py-1.5 rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50 text-xs"
					onclick={() => (showEditModal = false)}
				>
					Batal
				</button>
				<button
					type="submit"
					class="px-4 py-1.5 rounded-lg bg-indigo-600 hover:bg-indigo-700 text-white font-semibold text-xs shadow-2xs"
				>
					Perbarui Data Tenant
				</button>
			</div>
		</form>
	</Modal>
{/if}

<!-- Modal Konfirmasi Toggle Aktivasi / Non-Aktivasi -->
{#if targetCompany}
	<Modal bind:open={showToggleModal} size="sm" autoclose={false} title="Konfirmasi Perubahan Status">
		<div class="space-y-3 text-xs">
			{#if targetCompany.is_active}
				<div class="p-3 bg-rose-50 border border-rose-200 rounded-xl text-rose-800 space-y-1">
					<div class="font-bold flex items-center gap-1.5">
						<span>⚠️</span>
						<span>Peringatan Penonaktifan Tenant</span>
					</div>
					<p class="text-[11px] leading-relaxed">
						Anda akan menonaktifkan tenant <strong>{targetCompany.name}</strong>. Seluruh karyawan dan admin di perusahaan ini akan <strong>diblokir dari proses login dan presensi wajah/GPS</strong> hingga diaktifkan kembali.
					</p>
				</div>
			{:else}
				<div class="p-3 bg-emerald-50 border border-emerald-200 rounded-xl text-emerald-800 space-y-1">
					<div class="font-bold flex items-center gap-1.5">
						<span>✅</span>
						<span>Aktivasi Kembali Tenant</span>
					</div>
					<p class="text-[11px] leading-relaxed">
						Anda akan mengaktifkan kembali tenant <strong>{targetCompany.name}</strong>. Seluruh karyawan dan admin akan kembali dapat login dan melakukan presensi seperti biasa.
					</p>
				</div>
			{/if}

			<div class="flex justify-end gap-2 pt-2 border-t border-slate-100">
				<button
					type="button"
					class="px-3 py-1.5 rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50 text-xs"
					onclick={() => (showToggleModal = false)}
				>
					Batal
				</button>
				<button
					type="button"
					disabled={toggleLoading}
					onclick={confirmToggleStatus}
					class={`px-4 py-1.5 rounded-lg text-white font-semibold text-xs shadow-2xs ${
						targetCompany.is_active
							? 'bg-rose-600 hover:bg-rose-700'
							: 'bg-emerald-600 hover:bg-emerald-700'
					}`}
				>
					{toggleLoading ? 'Memproses...' : targetCompany.is_active ? 'Ya, Nonaktifkan Tenant' : 'Ya, Aktifkan Tenant'}
				</button>
			</div>
		</div>
	</Modal>
{/if}