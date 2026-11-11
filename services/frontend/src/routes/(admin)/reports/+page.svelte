<script lang="ts">
	import { onMount } from 'svelte';
	import { Alert, Spinner } from 'flowbite-svelte';
	import { auth } from '$lib/stores/auth.svelte';

	let loading = $state(true);
	let errorMsg = $state('');
	let reportData = $state<any>({ summary: {}, items: [] });

	// Filters
	let startDate = $state(new Date().toISOString().split('T')[0]);
	let endDate = $state(new Date().toISOString().split('T')[0]);
	let selectedStatus = $state('');

	const BACKEND_URL = 'http://localhost:8080/api/v1';

	onMount(loadReport);

	async function loadReport() {
		loading = true;
		errorMsg = '';
		try {
			let url = `${BACKEND_URL}/reports/attendance?start_date=${startDate}&end_date=${endDate}`;
			if (selectedStatus) {
				url += `&status=${selectedStatus}`;
			}

			const res = await fetch(url, {
				headers: { 'Authorization': `Bearer ${auth.token}` }
			});
			const json = await res.json();
			if (res.ok) {
				reportData = json.data || { summary: {}, items: [] };
			} else {
				errorMsg = json.error || 'Gagal mengambil laporan.';
			}
		} catch (err: any) {
			errorMsg = err.message;
		} finally {
			loading = false;
		}
	}
</script>

<div class="space-y-6">
	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pb-4 border-b border-slate-200">
		<div>
			<h1 class="text-2xl font-bold text-slate-900 tracking-tight">Laporan & Payroll Presensi</h1>
			<p class="text-xs text-slate-500 mt-0.5">Rekapitulasi kehadiran, keterlambatan, pulang awal, dan validasi biometrik</p>
		</div>
		<button
			type="button"
			onclick={() => window.print()}
			class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold text-xs rounded-lg shadow-sm shadow-indigo-600/20 transition-all"
		>
			<span>🖨️ Cetak / Ekspor Laporan</span>
		</button>
	</div>

	<!-- Filter Bar Card -->
	<div class="bg-white border border-slate-200/90 rounded-2xl p-5 shadow-2xs">
		<div class="grid grid-cols-1 sm:grid-cols-4 gap-3 items-end text-xs">
			<div>
				<label for="sdate" class="block mb-1 font-semibold text-slate-700">Tanggal Mulai</label>
				<input id="sdate" type="date" bind:value={startDate} class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs" />
			</div>

			<div>
				<label for="edate" class="block mb-1 font-semibold text-slate-700">Tanggal Akhir</label>
				<input id="edate" type="date" bind:value={endDate} class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs" />
			</div>

			<div>
				<label for="fstatus" class="block mb-1 font-semibold text-slate-700">Status Kehadiran</label>
				<select id="fstatus" bind:value={selectedStatus} class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-xs">
					<option value="">Semua Status</option>
					<option value="ON_TIME">Tepat Waktu (ON_TIME)</option>
					<option value="LATE">Terlambat (LATE)</option>
					<option value="OUT_OF_GEOFENCE">Luar Geofence</option>
				</select>
			</div>

			<div>
				<button 
					type="button" 
					onclick={loadReport}
					class="w-full py-2 bg-slate-900 hover:bg-slate-800 text-white font-semibold text-xs rounded-lg transition-colors"
				>
					Terapkan Filter
				</button>
			</div>
		</div>
	</div>

	{#if errorMsg}
		<Alert color="red" class="text-xs" dismissable>
			{errorMsg}
		</Alert>
	{/if}

	<!-- Summary Metrics -->
	<div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
		<div class="bg-white border border-slate-200/90 rounded-xl p-4 shadow-2xs">
			<span class="text-[11px] font-bold text-slate-400 uppercase">Total Log</span>
			<div class="text-2xl font-black text-slate-900 mt-1">{reportData.summary?.total_records || 0}</div>
		</div>
		<div class="bg-white border border-slate-200/90 rounded-xl p-4 shadow-2xs">
			<span class="text-[11px] font-bold text-slate-400 uppercase">Tepat Waktu</span>
			<div class="text-2xl font-black text-emerald-600 mt-1">{reportData.summary?.on_time || 0}</div>
		</div>
		<div class="bg-white border border-slate-200/90 rounded-xl p-4 shadow-2xs">
			<span class="text-[11px] font-bold text-slate-400 uppercase">Terlambat</span>
			<div class="text-2xl font-black text-amber-600 mt-1">{reportData.summary?.late || 0}</div>
		</div>
		<div class="bg-white border border-slate-200/90 rounded-xl p-4 shadow-2xs">
			<span class="text-[11px] font-bold text-slate-400 uppercase">Luar Geofence</span>
			<div class="text-2xl font-black text-rose-600 mt-1">{reportData.summary?.out_of_geofence || 0}</div>
		</div>
	</div>

	<!-- Report Table -->
	<div class="bg-white border border-slate-200/90 rounded-2xl p-6 shadow-2xs space-y-4">
		{#if loading}
			<div class="text-center py-12 text-slate-500 text-sm">
				<Spinner size="4" class="mr-2 inline" /> Memuat laporan presensi...
			</div>
		{:else if !reportData.records || reportData.records.length === 0}
			<div class="text-center py-12 text-slate-400 text-xs">
				Tidak ada catatan kehadiran pada rentang tanggal yang dipilih.
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full text-left text-xs">
					<thead>
						<tr class="text-[11px] font-bold text-slate-400 uppercase tracking-wider border-b border-slate-100">
							<th class="pb-3 pl-2">Tanggal</th>
							<th class="pb-3">Pegawai & NIK</th>
							<th class="pb-3">Cabang Kantor</th>
							<th class="pb-3">Waktu Masuk</th>
							<th class="pb-3">Status</th>
							<th class="pb-3">Biometrik Wajah</th>
							<th class="pb-3 pr-2 text-right">Jarak GPS</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100">
						{#each reportData.records as rec}
							<tr class="hover:bg-slate-50/70 transition-colors">
								<td class="py-3.5 pl-2 font-mono text-slate-700">{rec.attendance_date}</td>
								<td class="py-3.5 font-bold text-slate-900">{rec.employee_name}</td>
								<td class="py-3.5 text-slate-600">{rec.branch_name || 'Head Office'}</td>
								<td class="py-3.5 font-mono text-slate-800">
									{rec.check_in_time ? new Date(rec.check_in_time).toLocaleTimeString('id-ID') : '-'}
								</td>
								<td class="py-3.5">
									<span class={`px-2.5 py-0.5 rounded-full font-bold text-[10px] ${
										rec.check_in_status === 'ON_TIME'
											? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
											: rec.check_in_status === 'LATE'
											? 'bg-amber-50 text-amber-700 border border-amber-200'
											: 'bg-rose-50 text-rose-700 border border-rose-200'
									}`}>
										{rec.check_in_status}
									</span>
								</td>
								<td class="py-3.5 font-mono font-bold text-emerald-600">
									{Math.round((rec.check_in_similarity_score || 0) * 100)}% Match
								</td>
								<td class="py-3.5 pr-2 text-right font-mono text-slate-700">
									{Math.round(rec.check_in_distance_meters || 0)}m
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>