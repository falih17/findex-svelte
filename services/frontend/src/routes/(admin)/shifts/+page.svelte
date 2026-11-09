<script lang="ts">
	import { onMount } from 'svelte';
	import { 
		Card, Table, TableHead, TableHeadCell, TableBody, TableBodyRow, TableBodyCell, 
		Button, Badge, Modal, Label, Input, Alert, Spinner
	} from 'flowbite-svelte';
	import { auth } from '$lib/stores/auth.svelte';

	let shifts = $state<any[]>([]);
	let loading = $state(true);
	let errorMsg = $state('');
	let successMsg = $state('');

	// Form Tambah Shift
	let showModal = $state(false);
	let name = $state('Shift Reguler');
	let startTime = $state('08:00:00');
	let checkInDeadline = $state('08:30:00');
	let earliestCheckout = $state('16:30:00');
	let endTime = $state('17:00:00');

	const BACKEND_URL = 'http://localhost:8080/api/v1';

	onMount(loadShifts);

	async function loadShifts() {
		loading = true;
		errorMsg = '';
		try {
			const res = await fetch(`${BACKEND_URL}/schedules/shifts`, {
				headers: { 'Authorization': `Bearer ${auth.token}` }
			});
			const json = await res.json();
			if (res.ok) {
				shifts = json.data || [];
			} else {
				errorMsg = json.error || 'Gagal memuat jadwal shift.';
			}
		} catch (err: any) {
			errorMsg = err.message;
		} finally {
			loading = false;
		}
	}

	async function handleCreateShift(e: Event) {
		e.preventDefault();
		errorMsg = '';
		successMsg = '';

		try {
			const res = await fetch(`${BACKEND_URL}/schedules/shifts`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${auth.token}`
				},
				body: JSON.stringify({
					name,
					start_time: startTime,
					check_in_deadline: checkInDeadline,
					earliest_checkout_time: earliestCheckout,
					end_time: endTime,
					grace_period_late_minutes: 15,
					min_work_hours: 8
				})
			});

			const json = await res.json();
			if (!res.ok) throw new Error(json.error || 'Gagal membuat shift.');

			successMsg = `Jadwal shift "${name}" berhasil disimpan!`;
			showModal = false;
			loadShifts();
		} catch (err: any) {
			errorMsg = err.message;
		}
	}
</script>

<div class="space-y-6 max-w-5xl mx-auto">
	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 border-b border-gray-800 pb-4">
		<div>
			<h1 class="text-2xl font-bold text-white tracking-tight">Manajemen Jadwal Kerja & Shift</h1>
			<p class="text-xs text-gray-400">Atur jam mulai kerja, batas akhir check-in, dan waktu paling awal boleh check-out</p>
		</div>
		<Button color="green" size="sm" onclick={() => { showModal = true; errorMsg = ''; }}>
			➕ Tambah Shift Baru
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

	<!-- Shifts Table -->
	<Card class="bg-gray-900 border-gray-800 p-4 max-w-none shadow-xl">
		<Table hoverable class="text-xs text-left">
			<TableHead class="text-gray-400 border-b border-gray-800 bg-gray-800/40">
				<TableHeadCell>Nama Shift</TableHeadCell>
				<TableHeadCell>Jam Masuk Kerja</TableHeadCell>
				<TableHeadCell>Batas Akhir Check-In</TableHeadCell>
				<TableHeadCell>Awal Boleh Check-Out</TableHeadCell>
				<TableHeadCell>Jam Pulang Kerja</TableHeadCell>
				<TableHeadCell>Aturan Keterlambatan</TableHeadCell>
			</TableHead>
			<TableBody class="divide-y divide-gray-800">
				{#if loading}
					<TableBodyRow>
						<TableBodyCell colspan={6} class="text-center py-8 text-gray-500">
							<Spinner size="4" class="mr-2 inline" /> Memuat jadwal shift...
						</TableBodyCell>
					</TableBodyRow>
				{:else if shifts.length === 0}
					<TableBodyRow>
						<TableBodyCell colspan={6} class="text-center py-8 text-gray-500">
							Belum ada jadwal shift. Silakan buat shift baru.
						</TableBodyCell>
					</TableBodyRow>
				{:else}
					{#each shifts as s}
						<TableBodyRow class="hover:bg-gray-800/40">
							<TableBodyCell class="font-semibold text-white">{s.name}</TableBodyCell>
							<TableBodyCell class="font-mono text-emerald-400 font-bold">{s.start_time}</TableBodyCell>
							<TableBodyCell class="font-mono text-amber-400 font-semibold">
								{s.check_in_deadline || s.start_time}
								<span class="block text-[10px] text-gray-500 font-sans">Lewat batas = Late</span>
							</TableBodyCell>
							<TableBodyCell class="font-mono text-purple-400 font-semibold">
								{s.earliest_checkout_time || s.end_time}
								<span class="block text-[10px] text-gray-500 font-sans">Sebelum waktu ini = Early</span>
							</TableBodyCell>
							<TableBodyCell class="font-mono text-white">{s.end_time}</TableBodyCell>
							<TableBodyCell>
								<Badge color="dark" class="font-mono text-[10px]">Toleransi: {s.grace_period_late_minutes || 15}m</Badge>
							</TableBodyCell>
						</TableBodyRow>
					{/each}
				{/if}
			</TableBody>
		</Table>
	</Card>
</div>

<!-- Modal Tambah Shift -->
<Modal bind:open={showModal} size="md" autoclose={false} class="bg-gray-900 text-white">
	<h3 class="text-base font-bold text-white mb-4">Pengaturan Jadwal Shift Baru</h3>
	<form onsubmit={handleCreateShift} class="space-y-3.5 text-xs">
		<div>
			<Label class="mb-1 text-xs text-gray-300">Nama Shift</Label>
			<Input bind:value={name} required placeholder="cth: Shift Pagi Kantor / Shift Lapangan" class="bg-gray-800 text-xs" />
		</div>

		<div class="grid grid-cols-2 gap-2">
			<div>
				<Label class="mb-1 text-xs text-gray-300">Jam Mulai Masuk</Label>
				<Input type="time" step="1" bind:value={startTime} required class="bg-gray-800 text-xs" />
			</div>
			<div>
				<Label class="mb-1 text-xs text-amber-300 font-semibold">Batas Akhir Check-In</Label>
				<Input type="time" step="1" bind:value={checkInDeadline} required class="bg-gray-800 text-xs border-amber-500/50" />
				<span class="text-[10px] text-gray-400">Presensi setelah jam ini = Terlambat</span>
			</div>
		</div>

		<div class="grid grid-cols-2 gap-2">
			<div>
				<Label class="mb-1 text-xs text-purple-300 font-semibold">Waktu Awal Checkout</Label>
				<Input type="time" step="1" bind:value={earliestCheckout} required class="bg-gray-800 text-xs border-purple-500/50" />
				<span class="text-[10px] text-gray-400">Pulang sebelum ini = Early Leave</span>
			</div>
			<div>
				<Label class="mb-1 text-xs text-gray-300">Jam Selesai Kerja</Label>
				<Input type="time" step="1" bind:value={endTime} required class="bg-gray-800 text-xs" />
			</div>
		</div>

		<div class="flex justify-end gap-2 pt-2 border-t border-gray-800">
			<Button color="light" size="sm" onclick={() => showModal = false}>Batal</Button>
			<Button type="submit" color="green" size="sm" class="font-bold text-slate-950">Simpan Jadwal Shift</Button>
		</div>
	</form>
</Modal>
