<script lang="ts">
	import { onMount } from 'svelte';
	import { 
		Card, Table, TableHead, TableHeadCell, TableBody, TableBodyRow, TableBodyCell, 
		Badge, Button, Alert, Spinner, Indicator 
	} from 'flowbite-svelte';
	import { auth } from '$lib/stores/auth.svelte';

	let loading = $state(true);
	let history = $state<any[]>([]);
	let errorMsg = $state('');

	const BACKEND_URL = 'http://localhost:8080/api/v1';

	onMount(async () => {
		try {
			const empID = auth.user?.id || '55555555-5555-5555-5555-555555555555';
			const res = await fetch(`${BACKEND_URL}/attendance/history?employee_id=${empID}`, {
				headers: auth.token ? { 'Authorization': `Bearer ${auth.token}` } : {}
			});
			const json = await res.json();
			if (res.ok) {
				history = json.data || [];
			} else {
				errorMsg = json.error || 'Gagal memuat riwayat presensi.';
			}
		} catch (err: any) {
			errorMsg = err.message;
		} finally {
			loading = false;
		}
	});
</script>

<div class="max-w-4xl mx-auto space-y-5">
	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 border-b border-gray-800 pb-3">
		<div>
			<h1 class="text-xl font-bold text-white flex items-center gap-2">
				Riwayat Presensi Saya
				<Badge color="green" class="text-[10px] font-mono">LOG BIOMETRIK</Badge>
			</h1>
			<p class="text-xs text-gray-400 mt-0.5">Histori kehadiran, koordinat GPS & verifikasi kemiripan wajah</p>
		</div>
		<Button size="xs" color="dark" href="/" class="text-xs">
			← Kembali ke Kamera
		</Button>
	</div>

	{#if loading}
		<div class="flex justify-center items-center py-12 gap-3 text-gray-400 text-sm">
			<Spinner size="6" color="green" /> Memuat data histori presensi...
		</div>
	{:else if errorMsg}
		<Alert color="red" class="text-xs">
			{errorMsg}
		</Alert>
	{:else if history.length === 0}
		<Card padding="xl" class="bg-gray-900 border-gray-800 text-center max-w-none">
			<p class="text-sm font-semibold text-gray-300">Belum ada catatan presensi</p>
			<p class="text-xs text-gray-500 mt-1">Anda belum melakukan clock-in atau riwayat kosong.</p>
			<div class="mt-4">
				<Button size="xs" color="green" href="/">Mulai Presensi Sekarang</Button>
			</div>
		</Card>
	{:else}
		<Card padding="md" class="bg-gray-900 border-gray-800 max-w-none">
			<div class="overflow-x-auto">
				<Table hoverable={true} class="text-left text-xs text-gray-300">
					<TableHead class="bg-gray-800 text-gray-400 border-b border-gray-700">
						<TableHeadCell>Tanggal</TableHeadCell>
						<TableHeadCell>Waktu Masuk</TableHeadCell>
						<TableHeadCell>Jarak Kantor</TableHeadCell>
						<TableHeadCell>Biometrik Wajah</TableHeadCell>
						<TableHeadCell>Status Kehadiran</TableHeadCell>
						<TableHeadCell>Waktu Pulang</TableHeadCell>
					</TableHead>
					<TableBody class="divide-y divide-gray-800">
						{#each history as item}
							<TableBodyRow class="hover:bg-gray-800/50">
								<TableBodyCell class="font-semibold text-white">
									{item.attendance_date}
								</TableBodyCell>
								<TableBodyCell class="font-mono text-gray-300">
									{item.check_in_time ? new Date(item.check_in_time).toLocaleTimeString('id-ID') : '-'}
								</TableBodyCell>
								<TableBodyCell class="text-gray-300 font-mono">
									{Math.round(item.check_in_distance_meters || 0)}m
								</TableBodyCell>
								<TableBodyCell>
									<span class="font-mono font-bold text-emerald-400 bg-emerald-500/10 px-2 py-0.5 rounded border border-emerald-500/20">
										{Math.round((item.check_in_similarity_score || 0) * 100)}% Match
									</span>
								</TableBodyCell>
								<TableBodyCell>
									<Badge color={item.check_in_status === 'ON_TIME' ? 'green' : item.check_in_status === 'LATE' ? 'yellow' : 'red'} class="text-[10px] font-bold">
										{item.check_in_status}
									</Badge>
								</TableBodyCell>
								<TableBodyCell class="font-mono text-gray-400">
									{item.check_out_time ? new Date(item.check_out_time).toLocaleTimeString('id-ID') : '-'}
								</TableBodyCell>
							</TableBodyRow>
						{/each}
					</TableBody>
				</Table>
			</div>
		</Card>
	{/if}
</div>