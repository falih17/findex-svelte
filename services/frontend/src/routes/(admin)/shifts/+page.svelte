<script lang="ts">
	import { onMount } from 'svelte';
	import { 
		Modal, Label, Input, Alert, Spinner
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
					end_time: endTime
				})
			});

			const json = await res.json();
			if (!res.ok) {
				throw new Error(json.error || 'Gagal menyimpan shift.');
			}

			successMsg = `Shift "${name}" berhasil ditambahkan!`;
			showModal = false;
			loadShifts();
		} catch (err: any) {
			errorMsg = err.message;
		}
	}
</script>

<div class="space-y-6">
	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pb-4 border-b border-slate-200">
		<div>
			<h1 class="text-2xl font-bold text-slate-900 tracking-tight">Jadwal & Shift Kerja</h1>
			<p class="text-xs text-slate-500 mt-0.5">Konfigurasi jam masuk, batas telat, dan waktu awal checkout pegawai</p>
		</div>
		<button
			type="button"
			onclick={() => { showModal = true; errorMsg = ''; successMsg = ''; }}
			class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold text-xs rounded-lg shadow-sm shadow-indigo-600/20 transition-all"
		>
			<span>➕ Tambah Shift Baru</span>
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

	<!-- Shift Table Card -->
	<div class="bg-white border border-slate-200/90 rounded-2xl p-6 shadow-2xs space-y-4">
		{#if loading}
			<div class="text-center py-12 text-slate-500 text-sm">
				<Spinner size="4" class="mr-2 inline" /> Memuat daftar jadwal shift...
			</div>
		{:else if shifts.length === 0}
			<div class="text-center py-10 text-slate-400 text-xs">
				Belum ada jadwal shift. Klik "Tambah Shift Baru" untuk membuat jadwal kerja pertama.
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full text-left text-xs">
					<thead>
						<tr class="text-[11px] font-bold text-slate-400 uppercase tracking-wider border-b border-slate-100">
							<th class="pb-3 pl-2">Nama Shift</th>
							<th class="pb-3">Jam Masuk</th>
							<th class="pb-3">Batas Akhir Check-in</th>
							<th class="pb-3">Awal Pulang</th>
							<th class="pb-3">Jam Selesai</th>
							<th class="pb-3 pr-2 text-right">Status Toleransi</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100">
						{#each shifts as s}
							<tr class="hover:bg-slate-50/70 transition-colors">
								<td class="py-3.5 pl-2 font-bold text-slate-900">{s.name}</td>
								<td class="py-3.5 font-mono text-indigo-600 font-semibold">{s.start_time}</td>
								<td class="py-3.5 font-mono text-amber-600 font-semibold">
									<span class="px-2 py-0.5 rounded bg-amber-50 border border-amber-200">
										{s.check_in_deadline || s.start_time}
									</span>
								</td>
								<td class="py-3.5 font-mono text-slate-600">
									{s.earliest_checkout_time || s.end_time}
								</td>
								<td class="py-3.5 font-mono text-slate-900 font-bold">{s.end_time}</td>
								<td class="py-3.5 pr-2 text-right">
									<span class="px-2 py-0.5 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-700 font-bold text-[10px]">
										Aktif
									</span>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>

<!-- Modal Tambah Shift -->
<Modal bind:open={showModal} size="md" autoclose={false} title="Tambah Jadwal Shift Pegawai">
	<form onsubmit={handleCreateShift} class="space-y-4 text-xs">
		<div>
			<Label for="sname" class="mb-1 text-slate-700">Nama Shift Kerja</Label>
			<Input id="sname" bind:value={name} placeholder="Contoh: Shift Reguler Pagi" required />
		</div>

		<div class="grid grid-cols-2 gap-3">
			<div>
				<Label for="stime" class="mb-1 text-slate-700">Jam Mulai Masuk</Label>
				<Input id="stime" type="time" step="1" bind:value={startTime} required />
			</div>
			<div>
				<Label for="cdeadline" class="mb-1 text-amber-700 font-semibold">Batas Akhir Check-In</Label>
				<Input id="cdeadline" type="time" step="1" bind:value={checkInDeadline} required />
			</div>
		</div>

		<div class="grid grid-cols-2 gap-3">
			<div>
				<Label for="echeckout" class="mb-1 text-slate-700">Waktu Awal Checkout</Label>
				<Input id="echeckout" type="time" step="1" bind:value={earliestCheckout} required />
			</div>
			<div>
				<Label for="etime" class="mb-1 text-slate-700">Jam Selesai Shift</Label>
				<Input id="etime" type="time" step="1" bind:value={endTime} required />
			</div>
		</div>

		<div class="flex justify-end gap-2 pt-2 border-t border-slate-100">
			<button type="button" class="px-3 py-1.5 rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50 text-xs" onclick={() => showModal = false}>
				Batal
			</button>
			<button type="submit" class="px-4 py-1.5 rounded-lg bg-indigo-600 hover:bg-indigo-700 text-white font-semibold text-xs shadow-2xs">
				Simpan Jadwal Shift
			</button>
		</div>
	</form>
</Modal>