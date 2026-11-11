<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Modal, Label, Input, Select, Button, Alert } from 'flowbite-svelte';

	let { onToggleSidebar } = $props<{ onToggleSidebar?: () => void }>();

	let currentTime = $state('');
	let timer: any;

	// Manual Attendance Modal
	let showManualModal = $state(false);
	let manualEmpNik = $state('');
	let manualDate = $state(new Date().toISOString().split('T')[0]);
	let manualTime = $state('08:00');
	let manualType = $state('CHECK_IN');
	let manualNotes = $state('');
	let manualSuccess = $state('');

	onMount(() => {
		updateClock();
		timer = setInterval(updateClock, 1000);
	});

	onDestroy(() => {
		if (timer) clearInterval(timer);
	});

	function updateClock() {
		const now = new Date();
		const hours = String(now.getHours()).padStart(2, '0');
		const mins = String(now.getMinutes()).padStart(2, '0');
		const secs = String(now.getSeconds()).padStart(2, '0');
		currentTime = `${hours}:${mins}:${secs} WIB`;
	}

	function handleManualSubmit(e: Event) {
		e.preventDefault();
		manualSuccess = 'Presensi manual berhasil disimpan ke audit log!';
		setTimeout(() => {
			manualSuccess = '';
			showManualModal = false;
		}, 1200);
	}
</script>

<header class="sticky top-0 z-30 bg-white/95 backdrop-blur-md border-b border-slate-200/80 px-4 sm:px-6 py-3 flex items-center justify-between gap-3 shadow-2xs">
	<!-- Left: Mobile Hamburger & Search -->
	<div class="flex items-center gap-3 flex-1 min-w-0">
		<button
			type="button"
			class="lg:hidden p-2 rounded-lg text-slate-500 hover:text-slate-700 hover:bg-slate-100 transition-colors"
			onclick={onToggleSidebar}
			aria-label="Buka menu navigasi"
		>
			<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
			</svg>
		</button>

		<!-- Global Search Box -->
		<div class="relative w-full max-w-sm sm:max-w-md">
			<span class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-slate-400">
				<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
				</svg>
			</span>
			<input
				type="text"
				placeholder="Cari karyawan, NIK (mis: 2024018), atau divisi..."
				class="w-full pl-9 pr-4 py-2 bg-slate-100 hover:bg-slate-100/80 focus:bg-white text-xs sm:text-sm text-slate-800 placeholder-slate-400 rounded-lg border-0 focus:ring-2 focus:ring-indigo-500 transition-all"
			/>
		</div>
	</div>

	<!-- Right Actions & Clock -->
	<div class="flex items-center gap-2 sm:gap-3 shrink-0">
		<!-- Digital Clock WIB -->
		<div class="hidden md:flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-100 text-slate-700 font-mono text-xs font-semibold border border-slate-200/60">
			<svg class="w-3.5 h-3.5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
			</svg>
			<span>{currentTime || '08:00:00 WIB'}</span>
		</div>

		<!-- + Presensi Manual Button -->
		<button
			type="button"
			onclick={() => showManualModal = true}
			class="hidden sm:inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-slate-50 border border-slate-200 text-slate-700 font-semibold text-xs rounded-lg transition-all shadow-2xs hover:border-slate-300"
		>
			<svg class="w-3.5 h-3.5 text-slate-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
			</svg>
			<span>+ Presensi Manual</span>
		</button>

		<!-- Ekspor Rekap Button -->
		<button
			type="button"
			onclick={() => window.print()}
			class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-indigo-600 hover:bg-indigo-700 active:bg-indigo-800 text-white font-semibold text-xs rounded-lg transition-all shadow-sm shadow-indigo-600/20"
		>
			<svg class="w-3.5 h-3.5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
			</svg>
			<span>Ekspor Rekap</span>
		</button>

		<!-- Notification Bell -->
		<div class="relative">
			<button
				type="button"
				class="p-2 text-slate-500 hover:text-slate-800 hover:bg-slate-100 rounded-lg transition-colors relative"
				aria-label="Pemberitahuan"
			>
				<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
				</svg>
				<!-- Red Unread Badge Indicator -->
				<span class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-rose-500 ring-2 ring-white"></span>
			</button>
		</div>
	</div>
</header>

<!-- Modal Presensi Manual -->
<Modal bind:open={showManualModal} size="md" title="Formulir Presensi Manual HR">
	{#if manualSuccess}
		<Alert color="green" class="text-xs">{manualSuccess}</Alert>
	{/if}

	<form onsubmit={handleManualSubmit} class="space-y-3.5 text-xs">
		<p class="text-slate-500 text-[11px]">
			Gunakan formulir ini hanya untuk penyesuaian khusus (misalnya pegawai dinas luar, kendala GPS, atau surat tugas).
		</p>

		<div>
			<Label for="empNik" class="mb-1 text-slate-700">NIK / Nama Pegawai</Label>
			<Input id="empNik" bind:value={manualEmpNik} placeholder="Contoh: 2024091 atau Anindya Putri" required />
		</div>

		<div class="grid grid-cols-2 gap-3">
			<div>
				<Label for="manualDate" class="mb-1 text-slate-700">Tanggal</Label>
				<Input id="manualDate" type="date" bind:value={manualDate} required />
			</div>
			<div>
				<Label for="manualTime" class="mb-1 text-slate-700">Waktu (Jam:Menit)</Label>
				<Input id="manualTime" type="time" bind:value={manualTime} required />
			</div>
		</div>

		<div>
			<Label for="manualType" class="mb-1 text-slate-700">Jenis Presensi</Label>
			<Select id="manualType" bind:value={manualType}>
				<option value="CHECK_IN">Presensi Masuk (Clock-in)</option>
				<option value="CHECK_OUT">Presensi Pulang (Clock-out)</option>
			</Select>
		</div>

		<div>
			<Label for="manualNotes" class="mb-1 text-slate-700">Alasan / Keterangan Penyesuaian</Label>
			<Input id="manualNotes" bind:value={manualNotes} placeholder="Tugas luar kota / Gangguan jaringan" required />
		</div>

		<div class="flex justify-end gap-2 pt-2 border-t border-slate-100">
			<Button color="light" size="xs" onclick={() => showManualModal = false}>Batal</Button>
			<Button type="submit" color="purple" size="xs">Simpan Presensi Manual</Button>
		</div>
	</form>
</Modal>