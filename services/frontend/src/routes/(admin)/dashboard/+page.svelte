<script lang="ts">
	let stats = {
		totalEmployees: 48,
		presentToday: 42,
		onTime: 38,
		late: 4,
		attendanceRate: 87.5
	};

	let recentLogs = [
		{ name: 'Admin Findex', time: '07:54:12', status: 'ON_TIME', dist: '12m', similarity: '96%' },
		{ name: 'Budi Pratama', time: '08:02:40', status: 'ON_TIME', dist: '24m', similarity: '91%' },
		{ name: 'Siti Rahma', time: '08:16:05', status: 'LATE', dist: '35m', similarity: '89%' },
		{ name: 'Rian Nugraha', time: '08:22:19', status: 'OUT_OF_GEOFENCE', dist: '142m', similarity: '94%' },
	];
</script>

<div class="space-y-6">
	<!-- Page Header -->
	<header class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 border-b border-slate-800 pb-4">
		<div>
			<h1 class="text-2xl font-bold text-white tracking-tight">HR Executive Dashboard</h1>
			<p class="text-xs text-slate-400">Monitoring Kehadiran, Biometrik & Lokasi Pegawai Real-time</p>
		</div>
		<div class="flex items-center gap-2">
			<a href="/branches" class="px-3.5 py-1.5 bg-slate-800 hover:bg-slate-700 text-xs font-semibold rounded-lg border border-slate-700 text-slate-200">
				📍 Kelola Lokasi Geofence
			</a>
			<button class="px-3.5 py-1.5 bg-emerald-500 hover:bg-emerald-400 text-xs font-bold rounded-lg text-slate-950">
				📥 Export Laporan Excel
			</button>
		</div>
	</header>

	<!-- Metric Cards Grid -->
	<div class="grid grid-cols-2 sm:grid-cols-4 gap-3 sm:gap-4">
		<div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-1">
			<p class="text-xs text-slate-400">Total Pegawai Aktif</p>
			<p class="text-2xl font-black text-white">{stats.totalEmployees}</p>
			<p class="text-[10px] text-emerald-400">Semua Divisi & Cabang</p>
		</div>

		<div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-1">
			<p class="text-xs text-slate-400">Hadir Hari Ini</p>
			<p class="text-2xl font-black text-emerald-400">{stats.presentToday} <span class="text-xs font-normal text-slate-400">({stats.attendanceRate}%)</span></p>
			<p class="text-[10px] text-slate-400">Check-in tervalidasi</p>
		</div>

		<div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-1">
			<p class="text-xs text-slate-400">Tepat Waktu</p>
			<p class="text-2xl font-black text-white">{stats.onTime}</p>
			<p class="text-[10px] text-emerald-400">Sesuai jam shift</p>
		</div>

		<div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-1">
			<p class="text-xs text-slate-400">Anomali / Terlambat</p>
			<p class="text-2xl font-black text-amber-400">{stats.late}</p>
			<p class="text-[10px] text-amber-400/80">Perlu review HR</p>
		</div>
	</div>

	<!-- Live Ingestion Feed Table -->
	<div class="bg-slate-900 border border-slate-800 rounded-2xl p-5 space-y-4">
		<div class="flex items-center justify-between">
			<h2 class="text-sm font-bold text-white flex items-center gap-2">
				<span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
				Aktivitas Presensi Terkini Hari Ini
			</h2>
			<span class="text-[11px] text-slate-500 font-mono">Auto-refresh aktif</span>
		</div>

		<div class="overflow-x-auto">
			<table class="w-full text-left text-xs">
				<thead>
					<tr class="text-slate-500 border-b border-slate-800">
						<th class="pb-2 font-medium">Pegawai</th>
						<th class="pb-2 font-medium">Waktu Masuk</th>
						<th class="pb-2 font-medium">Jarak ke Kantor</th>
						<th class="pb-2 font-medium">Similarity Wajah</th>
						<th class="pb-2 font-medium">Status Geofence</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-800/60">
					{#each recentLogs as log}
						<tr>
							<td class="py-3 font-semibold text-white">{log.name}</td>
							<td class="py-3 font-mono text-slate-300">{log.time}</td>
							<td class="py-3 text-slate-300">{log.dist}</td>
							<td class="py-3 font-mono font-medium text-emerald-400">{log.similarity}</td>
							<td class="py-3">
								<span class="px-2 py-0.5 rounded text-[10px] font-bold uppercase {log.status === 'ON_TIME' ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : log.status === 'LATE' ? 'bg-amber-500/10 text-amber-400 border border-amber-500/20' : 'bg-red-500/10 text-red-400 border border-red-500/20'}">
									{log.status}
								</span>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
