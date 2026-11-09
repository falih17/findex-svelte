<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth.svelte';

	let loading = $state(true);
	let history = $state<any[]>([]);
	let errorMsg = $state('');

	const BACKEND_URL = 'http://localhost:8080/api/v1';

	onMount(async () => {
		try {
			const empID = auth.user?.id || '55555555-5555-5555-5555-555555555555';
			const res = await fetch(`${BACKEND_URL}/attendance/history?employee_id=${empID}`);
			const json = await res.json();
			if (res.ok) {
				history = json.data || [];
			} else {
				errorMsg = json.error || 'Gagal memuat riwayat.';
			}
		} catch (err: any) {
			errorMsg = err.message;
		} finally {
			loading = false;
		}
	});
</script>

<div class="max-w-2xl mx-auto space-y-4">
	<header class="flex items-center justify-between border-b border-slate-800 pb-3">
		<div>
			<h1 class="text-xl font-bold text-white">Riwayat Presensi</h1>
			<p class="text-xs text-slate-400">Histori kehadiran, jarak GPS & skor biometrik</p>
		</div>
		<a href="/" class="text-xs bg-slate-800 hover:bg-slate-700 text-slate-200 px-3 py-1.5 rounded-lg border border-slate-700">
			← Kembali ke Kamera
		</a>
	</header>

	{#if loading}
		<div class="text-center py-12 text-slate-500 text-sm">Memuat data histori presensi...</div>
	{:else if errorMsg}
		<div class="p-3 bg-red-500/10 border border-red-500/20 text-red-400 rounded-xl text-xs">
			{errorMsg}
		</div>
	{:else if history.length === 0}
		<div class="text-center py-12 bg-slate-900 border border-slate-800 rounded-2xl text-slate-500 text-sm">
			Belum ada catatan presensi untuk akun ini. Silakan coba presensi masuk di halaman utama!
		</div>
	{:else}
		<div class="space-y-2.5">
			{#each history as item}
				<div class="bg-slate-900 border border-slate-800 rounded-xl p-3.5 flex items-center justify-between hover:border-slate-700 transition-colors">
					<div class="space-y-1">
						<div class="flex items-center gap-2">
							<span class="font-semibold text-white text-sm">{item.attendance_date}</span>
							<span class="px-2 py-0.5 rounded text-[10px] font-bold tracking-wider uppercase {item.check_in_status === 'ON_TIME' ? 'bg-emerald-500/20 text-emerald-400' : 'bg-amber-500/20 text-amber-400'}">
								{item.check_in_status}
							</span>
						</div>
						<div class="text-[11px] text-slate-400">
							Masuk: <span class="text-slate-200 font-mono">{new Date(item.check_in_time).toLocaleTimeString('id-ID')}</span>
							• Jarak: <span class="text-slate-200">{Math.round(item.check_in_distance_meters || 0)}m</span>
						</div>
					</div>

					<div class="text-right">
						<div class="text-xs font-mono font-semibold text-emerald-400">
							{Math.round((item.check_in_similarity_score || 0) * 100)}% Match
						</div>
						<div class="text-[10px] text-slate-500">Biometrik Wajah</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
