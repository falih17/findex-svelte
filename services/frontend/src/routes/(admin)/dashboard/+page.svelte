<script lang="ts">
	import { onMount } from 'svelte';
	import { Modal, Button, Badge } from 'flowbite-svelte';
	import { auth } from '$lib/stores/auth.svelte';

	// KPI Stats
	let stats = $state({
		totalEmployees: 248,
		newEmployeesThisMonth: 3,
		presentToday: 218,
		onTimeRate: 87.9,
		lateCount: 18,
		lateToleranceRate: 7.2,
		avgLateMinutes: 14,
		permissionCount: 8,
		absentAlphaCount: 4
	});

	// Weekly compliance bar chart data (7 days)
	let weeklyTrend = [
		{ day: 'Senin', onTime: 220, late: 16, leave: 14 },
		{ day: 'Selasa', onTime: 212, late: 18, leave: 15 },
		{ day: 'Rabu', onTime: 232, late: 8, leave: 12 },
		{ day: 'Kamis', onTime: 226, late: 12, leave: 10 },
		{ day: 'Jumat', onTime: 208, late: 22, leave: 18 },
		{ day: 'Sabtu', onTime: 185, late: 10, leave: 6 },
		{ day: 'Hari Ini', onTime: 218, late: 18, leave: 12 }
	];

	// Approvals Queue
	let approvals = $state([
		{
			id: 1,
			name: 'Fajar Ramadhan',
			avatarBg: 'from-blue-500 to-indigo-600',
			deptRole: 'DevOps • Cuti Tahunan (2 Hari)',
			type: 'Cuti',
			typeColor: 'blue',
			reason: 'Keperluan keluarga mendesak di luar kota.',
			status: 'PENDING'
		},
		{
			id: 2,
			name: 'Dimas Wicaksono',
			avatarBg: 'from-amber-500 to-orange-600',
			deptRole: 'Warehouse • Koreksi Jam Absen',
			type: 'Koreksi',
			typeColor: 'amber',
			reason: 'GPS error saat tiba di Gerbang Muat pkl 07:55 WIB.',
			status: 'PENDING'
		}
	]);

	// Filter state for live logs table
	let activeFilter = $state<'ALL' | 'ON_TIME' | 'LATE' | 'REMOTE'>('ALL');

	// Real-time Attendance Logs
	let logs = $state([
		{
			id: 'att-1',
			name: 'Anindya Putri',
			nik: '2024091',
			avatarBg: 'from-rose-500 to-pink-600',
			jobTitle: 'UI/UX Designer',
			division: 'Product (08:30 - 17:30)',
			clockIn: '08:18 WIB',
			clockOut: '--:--',
			status: 'ON_TIME',
			statusLabel: 'Tepat Waktu',
			locationName: 'HQ Jakarta Pusat (Radius 8m)',
			biometricLabel: '99.4% Face Match',
			biometricType: 'MATCH',
			distanceMeters: 8,
			similarity: 0.994,
			photoUrl: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=120&auto=format&fit=crop&q=80'
		},
		{
			id: 'att-2',
			name: 'Budi Santoso',
			nik: '2023044',
			avatarBg: 'from-blue-600 to-cyan-600',
			jobTitle: 'Account Executive',
			division: 'Sales (09:00 - 18:00)',
			clockIn: '09:22 WIB',
			clockOut: '--:--',
			status: 'LATE',
			statusLabel: 'Telat 22 Menit',
			locationName: 'HQ Jakarta Pusat (Radius 15m)',
			biometricLabel: '98.1% Face Match',
			biometricType: 'MATCH',
			distanceMeters: 15,
			similarity: 0.981,
			photoUrl: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=120&auto=format&fit=crop&q=80'
		},
		{
			id: 'att-3',
			name: 'Rina Melati',
			nik: '2024102',
			avatarBg: 'from-emerald-600 to-teal-600',
			jobTitle: 'Senior Backend Engineer',
			division: 'Tech Remote (09:00 - 18:00)',
			clockIn: '08:55 WIB',
			clockOut: '--:--',
			status: 'REMOTE',
			statusLabel: 'WFH Disetujui',
			locationName: 'Bandung Wetan (Home Reg.)',
			biometricLabel: '99.8% Liveness Valid',
			biometricType: 'LIVENESS',
			distanceMeters: 0,
			similarity: 0.998,
			photoUrl: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=120&auto=format&fit=crop&q=80'
		},
		{
			id: 'att-4',
			name: 'Hendro Gunawan',
			nik: '2022081',
			avatarBg: 'from-slate-600 to-gray-700',
			jobTitle: 'Field Supervisor',
			division: 'Operations (08:00 - 17:00)',
			clockIn: '08:02 WIB',
			clockOut: '--:--',
			status: 'OUT_OF_GEOFENCE',
			statusLabel: 'Luar Geofence',
			locationName: 'Selisih +450m dari kantor',
			biometricLabel: 'Perlu Klarifikasi HR',
			biometricType: 'ANOMALY',
			distanceMeters: 450,
			similarity: 0.962,
			photoUrl: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=120&auto=format&fit=crop&q=80'
		}
	]);

	// Filtered logs computed
	let filteredLogs = $derived(
		activeFilter === 'ALL'
			? logs
			: activeFilter === 'ON_TIME'
			? logs.filter((l) => l.status === 'ON_TIME')
			: activeFilter === 'LATE'
			? logs.filter((l) => l.status === 'LATE')
			: logs.filter((l) => l.status === 'REMOTE' || l.status === 'OUT_OF_GEOFENCE')
	);

	// Detail Modal
	let showDetailModal = $state(false);
	let selectedLog = $state<any>(null);

	function openDetail(log: any) {
		selectedLog = log;
		showDetailModal = true;
	}

	function handleApprove(id: number) {
		approvals = approvals.filter((a) => a.id !== id);
	}

	function handleReject(id: number) {
		approvals = approvals.filter((a) => a.id !== id);
	}

	const BACKEND_URL = 'http://localhost:8080/api/v1';

	onMount(async () => {
		try {
			if (auth.token) {
				const res = await fetch(`${BACKEND_URL}/reports/attendance`, {
					headers: { Authorization: `Bearer ${auth.token}` }
				});
				if (res.ok) {
					const json = await res.json();
					if (json.data?.summary) {
						stats.totalEmployees = json.data.summary.total_records || 248;
						stats.presentToday = json.data.summary.total_records || 218;
						stats.onTimeRate = Math.round(
							((json.data.summary.on_time || 218) / (stats.presentToday || 1)) * 100
						);
					}
				}
			}
		} catch (err) {
			console.warn('Dashboard fetch fallback to preloaded state');
		}
	});
</script>

<div class="space-y-6">
	<!-- Top Page Header & Live Geofencing Status -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
		<div>
			<h1 class="text-2xl font-bold tracking-tight text-slate-900">
				Ringkasan Presensi Harian
			</h1>
			<p class="text-xs sm:text-sm text-slate-500 mt-0.5">
				Pantau absensi real-time, validasi geofence, dan tindakan persetujuan tim.
			</p>
		</div>

		<!-- Status Pills -->
		<div class="flex items-center gap-2.5">
			<!-- Sync Geofencing Aktif -->
			<div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-emerald-50 border border-emerald-300 text-emerald-700 text-xs font-semibold shadow-2xs">
				<span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
				<span>Sync Geofencing Aktif</span>
			</div>

			<!-- Tanggal Terkini -->
			<div class="inline-flex items-center px-3 py-1 rounded-lg bg-white border border-slate-200 text-slate-600 text-xs font-medium shadow-2xs">
				Senin, 09 November 2026
			</div>
		</div>
	</div>

	<!-- Metric KPI Cards (4 Grid Cards) -->
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
		<!-- Card 1: TOTAL KARYAWAN -->
		<div class="bg-white border border-slate-200/90 rounded-2xl p-5 shadow-2xs flex flex-col justify-between">
			<div class="flex items-start justify-between">
				<span class="text-[11px] font-bold tracking-wider text-slate-400 uppercase">
					Total Karyawan
				</span>
				<div class="w-8 h-8 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center">
					<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
					</svg>
				</div>
			</div>
			<div class="mt-2">
				<div class="text-3xl font-black text-slate-900 tracking-tight">
					{stats.totalEmployees} <span class="text-sm font-normal text-slate-500">Orang</span>
				</div>
				<div class="text-xs font-medium text-emerald-600 mt-1 flex items-center gap-1">
					<span>↗ +{stats.newEmployeesThisMonth} bergabung bulan ini</span>
				</div>
			</div>
		</div>

		<!-- Card 2: HADIR HARI INI -->
		<div class="bg-white border border-slate-200/90 rounded-2xl p-5 shadow-2xs flex flex-col justify-between">
			<div class="flex items-start justify-between">
				<span class="text-[11px] font-bold tracking-wider text-slate-400 uppercase">
					Hadir Hari Ini
				</span>
				<div class="w-8 h-8 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center">
					<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				</div>
			</div>
			<div class="mt-2">
				<div class="flex items-baseline gap-2">
					<span class="text-3xl font-black text-slate-900 tracking-tight">{stats.presentToday}</span>
					<span class="px-2 py-0.5 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-700 font-bold text-xs">
						{stats.onTimeRate}% On-Time
					</span>
				</div>
				<!-- Progress Bar -->
				<div class="w-full bg-slate-100 rounded-full h-1.5 mt-2.5 overflow-hidden">
					<div class="bg-emerald-500 h-1.5 rounded-full" style="width: {stats.onTimeRate}%;"></div>
				</div>
			</div>
		</div>

		<!-- Card 3: TERLAMBAT MASUK -->
		<div class="bg-white border border-slate-200/90 rounded-2xl p-5 shadow-2xs flex flex-col justify-between">
			<div class="flex items-start justify-between">
				<span class="text-[11px] font-bold tracking-wider text-slate-400 uppercase">
					Terlambat Masuk
				</span>
				<div class="w-8 h-8 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center">
					<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				</div>
			</div>
			<div class="mt-2">
				<div class="flex items-baseline gap-2">
					<span class="text-3xl font-black text-slate-900 tracking-tight">{stats.lateCount}</span>
					<span class="px-2 py-0.5 rounded-full bg-amber-50 border border-amber-200 text-amber-700 font-bold text-xs">
						{stats.lateToleranceRate}% Toleransi
					</span>
				</div>
				<div class="text-xs font-medium text-amber-600 mt-2 flex items-center gap-1">
					<span>⚠️ Rata-rata telat: {stats.avgLateMinutes} menit</span>
				</div>
			</div>
		</div>

		<!-- Card 4: TIDAK HADIR -->
		<div class="bg-white border border-slate-200/90 rounded-2xl p-5 shadow-2xs flex flex-col justify-between">
			<div class="flex items-start justify-between">
				<span class="text-[11px] font-bold tracking-wider text-slate-400 uppercase">
					Tidak Hadir
				</span>
				<div class="w-8 h-8 rounded-xl bg-rose-50 text-rose-500 flex items-center justify-center">
					<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7a4 4 0 11-8 0 4 4 0 018 0zM9 14a6 6 0 00-6 6v1h12v-1a6 6 0 00-6-6zM21 12h-6" />
					</svg>
				</div>
			</div>
			<div class="mt-2">
				<div class="grid grid-cols-2 gap-2 text-xs">
					<div>
						<div class="text-slate-400 font-medium">Izin/Cuti</div>
						<div class="text-2xl font-bold text-slate-800">{stats.permissionCount}</div>
					</div>
					<div>
						<div class="text-rose-500 font-medium">Alpha/Hilang</div>
						<div class="text-2xl font-bold text-rose-600">{stats.absentAlphaCount}</div>
					</div>
				</div>
				<div class="text-[11px] text-slate-400 mt-1.5">
					Perlu follow-up SPV divisi terkait
				</div>
			</div>
		</div>
	</div>

	<!-- Middle Section: Chart + Antrean Persetujuan -->
	<div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
		<!-- Left: Tren Kehadiran & Kepatuhan Shift (8 cols) -->
		<div class="lg:col-span-8 bg-white border border-slate-200/90 rounded-2xl p-6 shadow-2xs space-y-4">
			<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
				<div>
					<h2 class="text-base font-bold text-slate-900 flex items-center gap-2">
						<span class="text-indigo-600">📊</span>
						Tren Kehadiran & Kepatuhan Shift
					</h2>
					<p class="text-xs text-slate-500 mt-0.5">
						Perbandingan On-Time, Terlambat, dan Izin 7 hari terakhir.
					</p>
				</div>
				<select class="bg-slate-50 border border-slate-200 text-xs rounded-lg px-3 py-1.5 text-slate-700 font-medium focus:ring-2 focus:ring-indigo-500">
					<option>Semua Departemen</option>
					<option>Product & Engineering</option>
					<option>Sales & Marketing</option>
					<option>Operations & Field</option>
				</select>
			</div>

			<!-- Grouped Bar Chart Visualizer -->
			<div class="pt-4">
				<div class="h-64 flex items-end justify-between gap-2 sm:gap-4 px-2 border-b border-slate-200/80 pb-2">
					{#each weeklyTrend as item}
						{@const maxVal = 250}
						{@const onTimeH = (item.onTime / maxVal) * 100}
						{@const lateH = (item.late / maxVal) * 100}
						{@const leaveH = (item.leave / maxVal) * 100}

						<div class="flex-1 flex flex-col items-center gap-1 group relative">
							<!-- Hover Tooltip -->
							<div class="absolute -top-12 hidden group-hover:flex flex-col items-center bg-slate-900 text-white text-[10px] rounded px-2 py-1 z-20 whitespace-nowrap shadow-md">
								<span>{item.day}: {item.onTime} On-Time</span>
								<span class="text-amber-400">{item.late} Terlambat</span>
							</div>

							<!-- Bars Group -->
							<div class="w-full flex items-end justify-center gap-1 h-52">
								<!-- Tepat Waktu (Indigo) -->
								<div
									class="w-3 sm:w-4 rounded-t bg-indigo-600 hover:bg-indigo-500 transition-all"
									style="height: {onTimeH}%;"
								></div>
								<!-- Terlambat (Orange) -->
								<div
									class="w-3 sm:w-4 rounded-t bg-amber-500 hover:bg-amber-400 transition-all"
									style="height: {Math.max(lateH, 4)}%;"
								></div>
								<!-- Izin / Cuti (Slate/Light Gray) -->
								<div
									class="w-3 sm:w-4 rounded-t bg-slate-200 hover:bg-slate-300 transition-all"
									style="height: {Math.max(leaveH, 4)}%;"
								></div>
							</div>

							<!-- Day Label -->
							<span class="text-[11px] font-medium text-slate-500 mt-2">
								{item.day}
							</span>
						</div>
					{/each}
				</div>

				<!-- Legend -->
				<div class="flex items-center justify-center gap-6 mt-4 text-xs text-slate-600">
					<div class="flex items-center gap-2">
						<span class="w-3 h-3 rounded-xs bg-indigo-600"></span>
						<span>Tepat Waktu</span>
					</div>
					<div class="flex items-center gap-2">
						<span class="w-3 h-3 rounded-xs bg-amber-500"></span>
						<span>Terlambat</span>
					</div>
					<div class="flex items-center gap-2">
						<span class="w-3 h-3 rounded-xs bg-slate-200"></span>
						<span>Izin / Cuti</span>
					</div>
				</div>
			</div>
		</div>

		<!-- Right: Antrean Persetujuan (4 cols) -->
		<div class="lg:col-span-4 bg-white border border-slate-200/90 rounded-2xl p-6 shadow-2xs space-y-4 flex flex-col justify-between">
			<div>
				<div class="flex items-center justify-between pb-2 border-b border-slate-100">
					<div class="flex items-center gap-2">
						<span class="text-amber-500">📥</span>
						<h2 class="text-base font-bold text-slate-900">Antrean Persetujuan</h2>
					</div>
					<span class="px-2 py-0.5 rounded-full bg-amber-100 text-amber-800 text-[11px] font-bold">
						{approvals.length} Menunggu
					</span>
				</div>
				<p class="text-xs text-slate-500 mt-1">Permintaan cuti, koreksi & lembur</p>

				<!-- Approval Request Items -->
				<div class="mt-4 space-y-3.5">
					{#each approvals as req}
						<div class="border border-slate-200/90 rounded-xl p-3.5 bg-white space-y-2.5 hover:border-slate-300 transition-colors">
							<div class="flex items-start justify-between gap-2">
								<div class="flex items-center gap-2.5">
									<div class={`w-8 h-8 rounded-full bg-gradient-to-tr ${req.avatarBg} text-white font-bold text-xs flex items-center justify-center shrink-0`}>
										{req.name.slice(0, 2).toUpperCase()}
									</div>
									<div>
										<div class="text-xs font-bold text-slate-900">{req.name}</div>
										<div class="text-[10px] text-slate-500">{req.deptRole}</div>
									</div>
								</div>
								<span class={`text-[10px] font-bold px-2 py-0.5 rounded ${
									req.typeColor === 'blue'
										? 'bg-blue-50 text-blue-600 border border-blue-200'
										: 'bg-amber-50 text-amber-600 border border-amber-200'
								}`}>
									{req.type}
								</span>
							</div>

							<!-- Notes Quote Box -->
							<div class="bg-slate-50 rounded-lg p-2.5 text-xs text-slate-600 italic border border-slate-100">
								"{req.reason}"
							</div>

							<!-- Action Buttons -->
							<div class="flex items-center justify-end gap-2 pt-1">
								<button
									type="button"
									onclick={() => handleReject(req.id)}
									class="px-3 py-1.5 text-xs font-semibold text-slate-600 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition-colors"
								>
									Tolak
								</button>
								<button
									type="button"
									onclick={() => handleApprove(req.id)}
									class="px-3.5 py-1.5 text-xs font-semibold bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors shadow-2xs"
								>
									Setujui
								</button>
							</div>
						</div>
					{/each}

					{#if approvals.length === 0}
						<div class="text-center py-8 text-xs text-slate-400">
							🎉 Tidak ada antrean persetujuan baru.
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>

	<!-- Bottom Section: Log Presensi Real-Time Hari Ini -->
	<div class="bg-white border border-slate-200/90 rounded-2xl p-6 shadow-2xs space-y-4">
		<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pb-3 border-b border-slate-100">
			<div>
				<h2 class="text-base font-bold text-slate-900 flex items-center gap-2">
					<span class="text-indigo-600">🪪</span>
					Log Presensi Real-Time Hari Ini
				</h2>
				<p class="text-xs text-slate-500 mt-0.5">
					Menampilkan seluruh clock-in/out dengan status verifikasi biometrik & GPS.
				</p>
			</div>

			<!-- Filter Button Pills -->
			<div class="flex flex-wrap items-center gap-1.5">
				<button
					type="button"
					onclick={() => activeFilter = 'ALL'}
					class={`px-3 py-1.5 rounded-lg text-xs font-semibold transition-all ${
						activeFilter === 'ALL'
							? 'bg-slate-900 text-white'
							: 'bg-slate-100 text-slate-600 hover:bg-slate-200'
					}`}
				>
					Semua
				</button>
				<button
					type="button"
					onclick={() => activeFilter = 'ON_TIME'}
					class={`px-3 py-1.5 rounded-lg text-xs font-semibold transition-all ${
						activeFilter === 'ON_TIME'
							? 'bg-slate-900 text-white'
							: 'bg-slate-100 text-slate-600 hover:bg-slate-200'
					}`}
				>
					Tepat Waktu
				</button>
				<button
					type="button"
					onclick={() => activeFilter = 'LATE'}
					class={`px-3 py-1.5 rounded-lg text-xs font-semibold transition-all ${
						activeFilter === 'LATE'
							? 'bg-slate-900 text-white'
							: 'bg-slate-100 text-slate-600 hover:bg-slate-200'
					}`}
				>
					Terlambat
				</button>
				<button
					type="button"
					onclick={() => activeFilter = 'REMOTE'}
					class={`px-3 py-1.5 rounded-lg text-xs font-semibold transition-all ${
						activeFilter === 'REMOTE'
							? 'bg-slate-900 text-white'
							: 'bg-slate-100 text-slate-600 hover:bg-slate-200'
					}`}
				>
					WFH / Luar Radius
				</button>

				<!-- Refresh Button -->
				<button
					type="button"
					onclick={() => location.reload()}
					class="p-1.5 rounded-lg bg-slate-100 hover:bg-slate-200 text-slate-600 ml-1 transition-colors"
					title="Refresh data"
					aria-label="Refresh"
				>
					<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
					</svg>
				</button>
			</div>
		</div>

		<!-- Table -->
		<div class="overflow-x-auto">
			<table class="w-full text-left text-xs">
				<thead>
					<tr class="text-[11px] font-bold text-slate-400 uppercase tracking-wider border-b border-slate-100">
						<th class="pb-3 pl-2">Karyawan & NIK</th>
						<th class="pb-3">Shift & Divisi</th>
						<th class="pb-3">Clock In / Out</th>
						<th class="pb-3">Status Kehadiran</th>
						<th class="pb-3">Verifikasi & Lokasi</th>
						<th class="pb-3 pr-2 text-right">Aksi</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100">
					{#each filteredLogs as log}
						<tr class="hover:bg-slate-50/70 transition-colors">
							<!-- Karyawan & NIK -->
							<td class="py-3.5 pl-2">
								<div class="flex items-center gap-3">
									<div class={`w-9 h-9 rounded-full bg-gradient-to-tr ${log.avatarBg} text-white font-bold text-xs flex items-center justify-center shrink-0 shadow-2xs`}>
										{log.name.slice(0, 2).toUpperCase()}
									</div>
									<div>
										<div class="font-bold text-slate-900">{log.name}</div>
										<div class="text-[11px] text-slate-400 font-mono">NIK: {log.nik}</div>
									</div>
								</div>
							</td>

							<!-- Shift & Divisi -->
							<td class="py-3.5">
								<div class="font-medium text-slate-800">{log.jobTitle}</div>
								<div class="text-[11px] text-slate-500">{log.division}</div>
							</td>

							<!-- Clock In / Out -->
							<td class="py-3.5 font-mono">
								<div class="flex items-center gap-1.5 font-bold text-slate-900">
									<span class="w-2 h-2 rounded-full bg-emerald-500"></span>
									<span>{log.clockIn}</span>
								</div>
								<div class="text-[11px] text-slate-400">Out: {log.clockOut}</div>
							</td>

							<!-- Status Kehadiran -->
							<td class="py-3.5">
								{#if log.status === 'ON_TIME'}
									<span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-700 font-bold text-[11px]">
										✓ {log.statusLabel}
									</span>
								{:else if log.status === 'LATE'}
									<span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-amber-50 border border-amber-200 text-amber-700 font-bold text-[11px]">
										⏱ {log.statusLabel}
									</span>
								{:else if log.status === 'REMOTE'}
									<span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-blue-50 border border-blue-200 text-blue-700 font-bold text-[11px]">
										🏠 {log.statusLabel}
									</span>
								{:else}
									<span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-rose-50 border border-rose-200 text-rose-700 font-bold text-[11px]">
										🚫 {log.statusLabel}
									</span>
								{/if}
							</td>

							<!-- Verifikasi & Lokasi -->
							<td class="py-3.5">
								<div class="flex items-center gap-1 text-slate-700 font-medium">
									{#if log.status === 'OUT_OF_GEOFENCE'}
										<span class="text-rose-500">📍</span>
										<span class="text-rose-600 font-semibold">{log.locationName}</span>
									{:else}
										<span class="text-indigo-600">📍</span>
										<span>{log.locationName}</span>
									{/if}
								</div>
								<div class="text-[11px] mt-0.5">
									{#if log.status === 'OUT_OF_GEOFENCE'}
										<span class="text-amber-600 font-semibold flex items-center gap-1">
											⚠️ {log.biometricLabel}
										</span>
									{:else}
										<span class="text-emerald-600 font-semibold flex items-center gap-1">
											😊 {log.biometricLabel}
										</span>
									{/if}
								</div>
							</td>

							<!-- Aksi -->
							<td class="py-3.5 pr-2 text-right">
								<button
									type="button"
									onclick={() => openDetail(log)}
									class="inline-flex items-center gap-1 px-3 py-1 bg-white hover:bg-slate-50 border border-slate-200 text-slate-700 text-xs font-semibold rounded-lg shadow-2xs hover:border-slate-300 transition-colors"
								>
									<svg class="w-3.5 h-3.5 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
									</svg>
									<span>Detail</span>
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<!-- Pagination Footer -->
		<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pt-3 border-t border-slate-100 text-xs text-slate-500">
			<div>
				Menampilkan {filteredLogs.length} dari {stats.presentToday} log absensi terdata hari ini
			</div>
			<div class="flex items-center gap-1">
				<button class="px-2.5 py-1 rounded border border-slate-200 text-slate-600 hover:bg-slate-50 text-xs disabled:opacity-50" disabled>
					Sebelumnya
				</button>
				<button class="px-3 py-1 rounded bg-indigo-600 text-white font-bold text-xs">
					1
				</button>
				<button class="px-3 py-1 rounded border border-slate-200 text-slate-600 hover:bg-slate-50 text-xs">
					2
				</button>
				<button class="px-3 py-1 rounded border border-slate-200 text-slate-600 hover:bg-slate-50 text-xs">
					3
				</button>
				<button class="px-2.5 py-1 rounded border border-slate-200 text-slate-600 hover:bg-slate-50 text-xs">
					Selanjutnya
				</button>
			</div>
		</div>
	</div>
</div>

<!-- Modal Detail Absensi Pegawai -->
{#if selectedLog}
	<Modal bind:open={showDetailModal} size="md" title={`Detail Presensi - ${selectedLog.name}`}>
		<div class="space-y-4 text-xs">
			<div class="flex items-center gap-3 p-3 bg-slate-50 rounded-xl border border-slate-100">
				<div class={`w-12 h-12 rounded-xl bg-gradient-to-tr ${selectedLog.avatarBg} text-white font-bold text-base flex items-center justify-center shrink-0`}>
					{selectedLog.name.slice(0, 2).toUpperCase()}
				</div>
				<div>
					<div class="text-sm font-bold text-slate-900">{selectedLog.name}</div>
					<div class="text-slate-500">NIK: {selectedLog.nik} • {selectedLog.jobTitle}</div>
					<div class="text-indigo-600 font-medium text-[11px]">{selectedLog.division}</div>
				</div>
			</div>

			<div class="grid grid-cols-2 gap-3">
				<div class="p-3 bg-white border border-slate-200 rounded-xl space-y-1">
					<div class="text-[10px] uppercase font-bold text-slate-400">Waktu Masuk</div>
					<div class="text-base font-bold text-slate-900 font-mono">{selectedLog.clockIn}</div>
					<div class="text-[11px] text-slate-500">Pulang: {selectedLog.clockOut}</div>
				</div>

				<div class="p-3 bg-white border border-slate-200 rounded-xl space-y-1">
					<div class="text-[10px] uppercase font-bold text-slate-400">Verifikasi Wajah</div>
					<div class="text-base font-bold text-emerald-600 font-mono">
						{Math.round(selectedLog.similarity * 100)}% Match
					</div>
					<div class="text-[11px] text-emerald-600">Anti-Spoofing: Lulus (Live)</div>
				</div>
			</div>

			<div class="p-3 bg-white border border-slate-200 rounded-xl space-y-2">
				<div class="text-[10px] uppercase font-bold text-slate-400">Lokasi & Radius Geofence</div>
				<div class="flex items-center justify-between text-slate-700">
					<span>Titik Kantor:</span>
					<span class="font-medium">{selectedLog.locationName}</span>
				</div>
				<div class="flex items-center justify-between text-slate-700 border-t border-slate-100 pt-1.5">
					<span>Jarak Validasi GPS:</span>
					<span class="font-mono font-bold text-slate-900">{selectedLog.distanceMeters} meter</span>
				</div>
				<div class="flex items-center justify-between text-slate-700 border-t border-slate-100 pt-1.5">
					<span>Status Algoritma PostGIS:</span>
					{#if selectedLog.status === 'OUT_OF_GEOFENCE'}
						<span class="text-rose-600 font-bold">ST_DWithin FAILED (>100m)</span>
					{:else}
						<span class="text-emerald-600 font-bold">ST_DWithin PASSED (&le; 100m)</span>
					{/if}
				</div>
			</div>

			<div class="flex justify-end pt-2">
				<Button color="dark" size="xs" onclick={() => showDetailModal = false}>
					Tutup
				</Button>
			</div>
		</div>
	</Modal>
{/if}