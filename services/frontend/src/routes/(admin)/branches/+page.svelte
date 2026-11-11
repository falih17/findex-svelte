<script lang="ts">
	import { onMount } from 'svelte';
	import { 
		Modal, Label, Input, Textarea, Range, Alert, Spinner 
	} from 'flowbite-svelte';
	import { auth } from '$lib/stores/auth.svelte';

	let branches = $state<any[]>([]);
	let loading = $state(true);
	let errorMsg = $state('');
	let successMsg = $state('');

	// Form fields for adding new office location
	let showModal = $state(false);
	let branchName = $state('');
	let address = $state('');
	let latitude = $state(-6.2088);
	let longitude = $state(106.8456);
	let radiusMeters = $state(100);

	const BACKEND_URL = 'http://localhost:8080/api/v1';

	onMount(loadBranches);

	async function loadBranches() {
		loading = true;
		errorMsg = '';
		try {
			const res = await fetch(`${BACKEND_URL}/master/locations`, {
				headers: { 'Authorization': `Bearer ${auth.token}` }
			});
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
			const res = await fetch(`${BACKEND_URL}/master/locations`, {
				method: 'POST',
				headers: { 
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${auth.token}`
				},
				body: JSON.stringify({
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

			successMsg = `Lokasi kantor "${branchName}" berhasil disimpan!`;
			showModal = false;
			branchName = '';
			address = '';
			loadBranches();
		} catch (err: any) {
			errorMsg = err.message;
		}
	}
</script>

<div class="space-y-6">
	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pb-4 border-b border-slate-200">
		<div>
			<h1 class="text-2xl font-bold text-slate-900 tracking-tight">Geofencing & Lokasi Kantor</h1>
			<p class="text-xs text-slate-500 mt-0.5">Atur titik koordinat kantor dan batas radius presensi pegawai</p>
		</div>
		<button
			type="button"
			onclick={() => { showModal = true; errorMsg = ''; successMsg = ''; }}
			class="inline-flex items-center gap-1.5 px-3.5 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold text-xs rounded-lg shadow-sm shadow-indigo-600/20 transition-all"
		>
			<span>➕ Tambah Lokasi Kantor</span>
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

	<!-- Branches Cards Grid -->
	{#if loading}
		<div class="text-center py-12 text-slate-500 text-sm">
			<Spinner size="4" class="mr-2 inline" /> Memuat lokasi kantor...
		</div>
	{:else}
		<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
			{#each branches as branch}
				<div class="bg-white border border-slate-200/90 rounded-2xl p-5 space-y-3 shadow-2xs">
					<div class="flex items-start justify-between">
						<div>
							<h3 class="font-bold text-slate-900 text-base">{branch.name}</h3>
							<p class="text-xs text-slate-500 mt-0.5">{branch.address || 'Alamat belum diatur'}</p>
						</div>
						<span class="px-2.5 py-1 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-700 font-mono font-bold text-xs">
							Radius: {branch.radius_meters}m
						</span>
					</div>

					<div class="bg-slate-50 rounded-xl p-3 text-xs space-y-1 font-mono text-slate-600 border border-slate-100">
						<div>Latitude: <span class="text-slate-900 font-bold">{branch.latitude}</span></div>
						<div>Longitude: <span class="text-slate-900 font-bold">{branch.longitude}</span></div>
					</div>

					<div class="text-[11px] text-slate-400 flex items-center justify-between pt-2 border-t border-slate-100">
						<span>Validasi: PostGIS Spatial SRID 4326</span>
						<span class="text-emerald-600 font-semibold">✓ Geofence Aktif</span>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Modal Tambah Lokasi -->
<Modal bind:open={showModal} size="md" autoclose={false} title="Tambah Titik Lokasi Kantor & Geofence">
	<form onsubmit={handleCreateBranch} class="space-y-4 text-xs">
		<div>
			<Label for="bname" class="mb-1 text-slate-700">Nama Kantor / Cabang</Label>
			<Input id="bname" bind:value={branchName} placeholder="Contoh: Kantor Pusat (HQ) Jakarta" required />
		</div>

		<div>
			<Label for="baddr" class="mb-1 text-slate-700">Alamat Lengkap</Label>
			<Textarea id="baddr" bind:value={address} rows={2} placeholder="Jl. Jendral Sudirman No. 1..." required />
		</div>

		<div class="grid grid-cols-2 gap-3">
			<div>
				<Label for="lat" class="mb-1 text-slate-700">Latitude</Label>
				<Input id="lat" type="number" step="0.000001" bind:value={latitude} required />
			</div>
			<div>
				<Label for="lng" class="mb-1 text-slate-700">Longitude</Label>
				<Input id="lng" type="number" step="0.000001" bind:value={longitude} required />
			</div>
		</div>

		<button 
			type="button" 
			onclick={detectCurrentGPS} 
			class="w-full py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-lg text-xs font-semibold flex items-center justify-center gap-2 border border-slate-200 transition-colors"
		>
			📍 Deteksi Koordinat GPS Saya Sekarang
		</button>

		<div class="space-y-1">
			<div class="flex justify-between items-center text-slate-700">
				<Label for="radius">Radius Geofence Presensi</Label>
				<span class="font-mono font-bold text-indigo-600">{radiusMeters} meter</span>
			</div>
			<Range id="radius" min="20" max="1000" step="10" bind:value={radiusMeters} class="w-full" />
			<p class="text-[10px] text-slate-400">Pegawai di luar radius ini akan ditandai OUT_OF_GEOFENCE.</p>
		</div>

		<div class="flex justify-end gap-2 pt-2 border-t border-slate-100">
			<button type="button" class="px-3 py-1.5 rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50 text-xs" onclick={() => showModal = false}>
				Batal
			</button>
			<button type="submit" class="px-4 py-1.5 rounded-lg bg-indigo-600 hover:bg-indigo-700 text-white font-semibold text-xs shadow-2xs">
				Simpan Lokasi
			</button>
		</div>
	</form>
</Modal>