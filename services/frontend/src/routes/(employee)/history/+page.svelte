<script lang="ts">
	import { onMount } from 'svelte';
	import { Alert, Spinner } from 'flowbite-svelte';
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

<div class="space-y-6">
	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pb-4 border-b border-slate-200">
		<div>
			<h1 class="text-2xl font-bold text-slate-900 tracking-tight">Riwayat Presensi Saya</h1>
			<p class="text-xs text-slate-500 mt-0.5">Histori kehadiran, koordinat GPS & verifikasi kemiripan wajah biometrik</p>
		</div>
		<a
			href="/"
			class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-white hover:bg-slate-50 border border-slate-200 text-slate-700 font-semibold text-xs rounded-lg shadow-2xs transition-all"
		>
			<span>← Kembali ke Kamera</span>
		</a>
	</div>

	{#if loading}
		<div class="text-center py-12 text-slate-500 text-sm">
			<Spinner size="4" class="mr-2 inline" /> Memuat data histori presensi...
		</div>
	{:else if errorMsg}
		<Alert color="red" class="text-xs">
			{errorMsg}
		</Alert>
	{:else if history.length === 0}
		<div class="bg-white border border-slate-200/90 rounded-2xl p-12 text-center shadow-2xs space-y-3">
			<div class="w-12 h-12 rounded-xl bg-slate-100 text-slate-400 flex items-center justify-center mx-auto text-xl">
				📅
			</div>
			<div class="text-sm font-bold text-slate-800">Belum Ada Catatan Presensi</div>
			<p class="text-xs text-slate-400 max-w-sm mx-auto">
				Anda belum melakukan clock-in untuk akun ini. Silakan mulai presensi di halaman utama.
			</p>
			<div class="pt-2">
				<a href="/" class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-xs font-semibold shadow-2xs inline-block">
					Presensi Sekarang
				</a>
			</div>
		</div>
	{:else}
		<div class="bg-white border border-slate-200/90 rounded-2xl p-6 shadow-2xs space-y-4">
			<div class="overflow-x-auto">
				<table class="w-full text-left text-xs">
					<thead>
						<tr class="text-[11px] font-bold text-slate-400 uppercase tracking-wider border-b border-slate-100">
							<th class="pb-3 pl-2">Tanggal</th>
							<th class="pb-3">Waktu Masuk</th>
							<th class="pb-3">Jarak Kantor</th>
							<th class="pb-3">Biometrik Wajah</th>
							<th class="pb-3">Status Kehadiran</th>
							<th class="pb-3 pr-2 text-right">Waktu Pulang</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100">
						{#each history as item}
							<tr class="hover:bg-slate-50/70 transition-colors">
								<td class="py-3.5 pl-2 font-bold text-slate-900">{item.attendance_date}</td>
								<td class="py-3.5 font-mono text-slate-800">
									{item.check_in_time ? new Date(item.check_in_time).toLocaleTimeString('id-ID') : '-'}
								</td>
								<td class="py-3.5 font-mono text-slate-600">
									{Math.round(item.check_in_distance_meters || 0)}m
								</td>
								<td class="py-3.5 font-mono font-bold text-emerald-600">
									{Math.round((item.check_in_similarity_score || 0) * 100)}% Match
								</td>
								<td class="py-3.5">
									<span class={`px-2.5 py-0.5 rounded-full font-bold text-[10px] ${
										item.check_in_status === 'ON_TIME'
											? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
											: 'bg-amber-50 text-amber-700 border border-amber-200'
									}`}>
										{item.check_in_status}
									</span>
								</td>
								<td class="py-3.5 pr-2 text-right font-mono text-slate-500">
									{item.check_out_time ? new Date(item.check_out_time).toLocaleTimeString('id-ID') : '-'}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	{/if}
</div>