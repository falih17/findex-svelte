<script lang="ts">
	import { onMount } from 'svelte';

	let branches = $state<any[]>([]);
	let loading = $state(true);
	let errorMsg = $state('');
	let successMsg = $state('');

	// Form fields for adding new office
	let showModal = $state(false);
	let branchName = $state('');
	let address = $state('');
	let latitude = $state(-6.2088);
	let longitude = $state(106.8456);
	let radiusMeters = $state(100);

	const BACKEND_URL = 'http://localhost:8080/api/v1';
	const COMPANY_ID = '11111111-1111-1111-1111-111111111111';

	onMount(loadBranches);

	async function loadBranches() {
		loading = true;
		errorMsg = '';
		try {
			const res = await fetch(`${BACKEND_URL}/branches?company_id=${COMPANY_ID}`);
			const json = await res.json();
			if (res.ok) {
				branches = json.data || [];
			} else {
				errorMsg = json.error || 'Gagal memuat cabang.';
			}
		} catch (err: any) {
			errorMsg = err.message;
		} finally {
			loading = false;
		}
	}

	function detectCurrentGPS() {
		if ('geolocation' in navigator) {
			navigator.geolocation.getCurrentPosition((pos) => {
				latitude = Number(pos.coords.latitude.toFixed(6));
				longitude = Number(pos.coords.longitude.toFixed(6));
			});
		}
	}

	async function handleCreateBranch(e: Event) {
		e.preventDefault();
		errorMsg = '';
		successMsg = '';

		try {
			const res = await fetch(`${BACKEND_URL}/branches`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					company_id: COMPANY_ID,
					name: branchName,
					address: address,
					latitude: Number(latitude),
					longitude: Number(longitude),
					radius_meters: Number(radiusMeters)
				})
			});

			const json = await res.json();
			if (!res.ok) {
				throw new Error(json.error || 'Gagal menyimpan lokasi kantor.');
			}

			successMsg = `Lokasi kantor "${branchName}" berhasil ditambahkan!`;
			showModal = false;
			branchName = '';
			address = '';
			loadBranches();
		} catch (err: any) {
			errorMsg = err.message;
		}
	}
</script>

<div class="space-y-6 max-w-4xl mx-auto">
	<!-- Page Header -->
	<header class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 border-b border-slate-800 pb-4">
		<div>
			<h1 class="text-2xl font-bold text-white tracking-tight">Pengaturan Lokasi Absen (Geofencing)</h1>
			<p class="text-xs text-slate-400">Tentukan kantor pusat, cabang, atau lokasi proyek beserta radius toleransi presensi</p>
		</div>
		<button
			onclick={() => { showModal = true; errorMsg = ''; successMsg = ''; }}
			class="px-4 py-2 bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold text-xs rounded-xl shadow-lg shadow-emerald-500/20 transition-all flex items-center gap-1.5 self-start"
		>
			<span>➕ Tambah Lokasi Kantor</span>
		</button>
	</header>

	{#if successMsg}
		<div class="p-3 bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 rounded-xl text-xs">
			✅ {successMsg}
		</div>
	{/if}

	{#if errorMsg}
		<div class="p-3 bg-red-500/10 border border-red-500/20 text-red-400 rounded-xl text-xs">
			⚠️ {errorMsg}
		</div>
	{/if}

	<!-- Branches Cards List -->
	{#if loading}
		<div class="text-center py-12 text-slate-500 text-sm">Memuat lokasi kantor...</div>
	{:else}
		<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
			{#each branches as branch}
				<div class="bg-slate-900 border border-slate-800 rounded-2xl p-5 space-y-3 hover:border-slate-700 transition-all">
					<div class="flex items-start justify-between">
						<div>
							<h3 class="font-bold text-white text-base">{branch.name}</h3>
							<p class="text-xs text-slate-400">{branch.address || 'Alamat belum diatur'}</p>
						</div>
						<span class="px-2.5 py-1 rounded-full bg-emerald-500/10 text-emerald-400 text-[11px] font-semibold border border-emerald-500/20">
							Radius: {branch.radius_meters}m
						</span>
					</div>

					<div class="bg-slate-950 rounded-xl p-3 text-xs space-y-1 font-mono text-slate-400 border border-slate-800/80">
						<div>Latitude: <span class="text-slate-200">{branch.latitude}</span></div>
						<div>Longitude: <span class="text-slate-200">{branch.longitude}</span></div>
					</div>

					<div class="text-[11px] text-slate-500 flex items-center justify-between pt-1">
						<span>Sistem Geospasial: PostGIS WGS84</span>
						<span class="text-emerald-400 font-medium">Aktif untuk Presensi</span>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Modal Tambah Lokasi -->
{#if showModal}
	<div class="fixed inset-0 bg-black/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
		<div class="bg-slate-900 border border-slate-800 rounded-2xl max-w-lg w-full p-6 space-y-4 shadow-2xl">
			<div class="flex items-center justify-between border-b border-slate-800 pb-3">
				<h2 class="text-base font-bold text-white">Tambah Lokasi Kantor / Geofence</h2>
				<button onclick={() => showModal = false} class="text-slate-400 hover:text-white text-sm">✕</button>
			</div>

			<form onsubmit={handleCreateBranch} class="space-y-3.5 text-xs">
				<div>
					<label class="block font-medium text-slate-300 mb-1">Nama Lokasi / Cabang</label>
					<input
						type="text"
						bind:value={branchName}
						required
						placeholder="cth: Kantor Cabang Surabaya / Proyek MRT"
						class="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-emerald-500"
					/>
				</div>

				<div>
					<label class="block font-medium text-slate-300 mb-1">Alamat Lengkap</label>
					<textarea
						bind:value={address}
						rows="2"
						placeholder="cth: Jl. Pemuda No. 45, Genteng, Surabaya"
						class="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-emerald-500"
					></textarea>
				</div>

				<div class="grid grid-cols-2 gap-2">
					<div>
						<label class="block font-medium text-slate-300 mb-1">Latitude</label>
						<input
							type="number"
							step="any"
							bind:value={latitude}
							required
							class="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-emerald-500"
						/>
					</div>
					<div>
						<label class="block font-medium text-slate-300 mb-1">Longitude</label>
						<input
							type="number"
							step="any"
							bind:value={longitude}
							required
							class="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-emerald-500"
						/>
					</div>
				</div>

				<button
					type="button"
					onclick={detectCurrentGPS}
					class="w-full py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium rounded-lg text-[11px] border border-slate-700"
				>
					📍 Gunakan Titik Koordinat GPS Saat Ini
				</button>

				<div>
					<div class="flex items-center justify-between mb-1">
						<label class="font-medium text-slate-300">Radius Toleransi Presensi (Meter):</label>
						<span class="font-mono text-emerald-400 font-bold">{radiusMeters} meter</span>
					</div>
					<input
						type="range"
						min="10"
						max="500"
						step="5"
						bind:value={radiusMeters}
						class="w-full accent-emerald-500 cursor-pointer"
					/>
					<p class="text-[10px] text-slate-500 mt-1">Pegawai hanya dapat melakukan presensi jika berada dalam batas radius ini.</p>
				</div>

				<div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-800">
					<button
						type="button"
						onclick={() => showModal = false}
						class="px-4 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-xl text-xs"
					>
						Batal
					</button>
					<button
						type="submit"
						class="px-4 py-2 bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold rounded-xl text-xs shadow-md shadow-emerald-500/20"
					>
						Simpan Lokasi Kantor
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
