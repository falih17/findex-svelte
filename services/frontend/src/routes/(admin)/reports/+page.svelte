<script lang="ts">
	import { onMount } from 'svelte';
	import { 
		Card, Table, TableHead, TableHeadCell, TableBody, TableBodyRow, TableBodyCell, 
		Button, Badge, Input, Select, Alert, Spinner, Label
	} from 'flowbite-svelte';
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

<div class="space-y-6 max-w-5xl mx-auto">
	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 border-b border-gray-800 pb-4">
		<div>
			<h1 class="text-2xl font-bold text-white tracking-tight">Kelompok Laporan Presensi Pegawai</h1>
			<p class="text-xs text-gray-400">Rekapitulasi kehadiran, keterlambatan, pulang awal, dan validasi biometrik</p>
		</div>
		<Button color="green" size="sm" onclick={() => window.print()}>
			🖨️ Cetak / Ekspor Laporan
		</Button>
	</div>

	<!-- Filter Bar Card -->
	<Card class="bg-gray-900 border-gray-800 p-4 max-w-none shadow-xl">
		<div class="grid grid-cols-1 sm:grid-cols-4 gap-3 text-xs items-end">
			<div>
				<Label class="mb-1 text-xs text-gray-300">Dari Tanggal</Label>
				<Input type="date" bind:value={startDate} class="bg-gray-800 text-xs" />
			</div>
			<div>
				<Label class="mb-1 text-xs text-gray-300">Sampai Tanggal</Label>
				<Input type="date" bind:value={endDate} class="bg-gray-800 text-xs" />
			</div>
			<div>
				<Label class="mb-1 text-xs text-gray-300">Status Kehadiran</Label>
				<Select bind:value={selectedStatus} class="bg-gray-800 text-xs">
					<option value="">Semua Status</option>
					<option value="ON_TIME">Tepat Waktu (ON_TIME)</option>
					<option value="LATE">Terlambat (LATE)</option>
					<option value="OUT_OF_GEOFENCE">Di Luar Lokasi (OUT_OF_GEOFENCE)</option>
					<option value="WFH">Work From Home (WFH)</option>
				</Select>
			</div>
			<Button color="blue" size="sm" class="w-full text-xs" onclick={loadReport}>
				🔍 Terapkan Filter
			</Button>
		</div>
	</Card>

	<!-- Metric Summary Cards -->
	{#if reportData.summary}
		<div class="grid grid-cols-2 sm:grid-cols-5 gap-3">
			<Card class="bg-gray-900 border-gray-800 p-3 space-y-1">
				<p class="text-[11px] text-gray-400">Total Presensi</p>
				<p class="text-xl font-bold text-white">{reportData.summary.total_records || 0}</p>
			</Card>
			<Card class="bg-gray-900 border-gray-800 p-3 space-y-1">
				<p class="text-[11px] text-gray-400">Tepat Waktu</p>
				<p class="text-xl font-bold text-emerald-400">{reportData.summary.on_time_count || 0}</p>
			</Card>
			<Card class="bg-gray-900 border-gray-800 p-3 space-y-1">
				<p class="text-[11px] text-gray-400">Terlambat</p>
				<p class="text-xl font-bold text-amber-400">{reportData.summary.late_count || 0}</p>
			</Card>
			<Card class="bg-gray-900 border-gray-800 p-3 space-y-1">
				<p class="text-[11px] text-gray-400">Pulang Awal</p>
				<p class="text-xl font-bold text-purple-400">{reportData.summary.early_leave_count || 0}</p>
			</Card>
			<Card class="bg-gray-900 border-gray-800 p-3 space-y-1">
				<p class="text-[11px] text-gray-400">Rata-rata Wajah</p>
				<p class="text-xl font-bold text-blue-400">
					{Math.round((reportData.summary.avg_face_similarity || 0) * 100)}%
				</p>
			</Card>
		</div>
	{/if}

	{#if errorMsg}
		<Alert color="red" class="text-xs" dismissable>{errorMsg}</Alert>
	{/if}

	<!-- Report Table -->
	<Card class="bg-gray-900 border-gray-800 p-4 max-w-none shadow-xl">
		<Table hoverable class="text-xs text-left">
			<TableHead class="text-gray-400 border-b border-gray-800 bg-gray-800/40">
				<TableHeadCell>Tanggal</TableHeadCell>
				<TableHeadCell>Pegawai</TableHeadCell>
				<TableHeadCell>Lokasi Kantor</TableHeadCell>
				<TableHeadCell>Jam Masuk</TableHeadCell>
				<TableHeadCell>Jam Pulang</TableHeadCell>
				<TableHeadCell>Status & Jarak</TableHeadCell>
				<TableHeadCell>Skor Biometrik</TableHeadCell>
			</TableHead>
			<TableBody class="divide-y divide-gray-800">
				{#if loading}
					<TableBodyRow>
						<TableBodyCell colspan={7} class="text-center py-8 text-gray-500">
							<Spinner size="4" class="mr-2 inline" /> Menghasilkan laporan presensi...
						</TableBodyCell>
					</TableBodyRow>
				{:else if reportData.items.length === 0}
					<TableBodyRow>
						<TableBodyCell colspan={7} class="text-center py-8 text-gray-500">
							Tidak ada data presensi pada rentang tanggal atau filter yang dipilih.
						</TableBodyCell>
					</TableBodyRow>
				{:else}
					{#each reportData.items as item}
						<TableBodyRow class="hover:bg-gray-800/40">
							<TableBodyCell class="font-semibold text-white">{item.attendance_date}</TableBodyCell>
							<TableBodyCell>
								<div class="font-semibold text-white">{item.employee_name}</div>
								<div class="text-[10px] text-gray-400 font-mono">{item.employee_nik}</div>
							</TableBodyCell>
							<TableBodyCell class="text-gray-300">{item.branch_name}</TableBodyCell>
							<TableBodyCell class="font-mono text-emerald-400">{item.check_in_time || '-'}</TableBodyCell>
							<TableBodyCell class="font-mono text-purple-400">{item.check_out_time || '-'}</TableBodyCell>
							<TableBodyCell>
								<Badge color={item.check_in_status === 'ON_TIME' ? 'green' : item.check_in_status === 'LATE' ? 'yellow' : 'red'} class="text-[10px]">
									{item.check_in_status}
								</Badge>
								<span class="block text-[10px] text-gray-500 font-mono mt-0.5">{Math.round(item.check_in_distance_meters || 0)}m</span>
							</TableBodyCell>
							<TableBodyCell class="font-mono text-emerald-400 font-bold">
								{Math.round((item.face_similarity_score || 0) * 100)}%
							</TableBodyCell>
						</TableBodyRow>
					{/each}
				{/if}
			</TableBody>
		</Table>
	</Card>
</div>
