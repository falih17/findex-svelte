<script lang="ts">
	import { onMount } from 'svelte';
	import { 
		Card, Button, Badge, Modal, Label, Input, Textarea, Range, Alert, Spinner 
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

<div class="space-y-6 max-w-5xl mx-auto">
	<!-- Page Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 border-b border-gray-800 pb-4">
		<div>
			<h1 class="text-2xl font-bold text-white tracking-tight">Manajemen Lokasi Absen (Geofencing)</h1>
			<p class="text-xs text-gray-400">Atur titik lokasi kantor dan batas radius di mana pegawai dapat melakukan presensi</p>
		</div>
		<Button color="green" size="sm" onclick={() => { showModal = true; errorMsg = ''; successMsg = ''; }}>
			➕ Tambah Lokasi Kantor
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

	<!-- Branches Cards Grid -->
	{#if loading}
		<div class="text-center py-12 text-gray-500 text-sm">
			<Spinner size="4" class="mr-2 inline" /> Memuat lokasi kantor...
		</div>
	{:else}
		<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
			{#each branches as branch}
				<Card class="bg-gray-900 border-gray-800 p-5 space-y-3 shadow-xl max-w-none">
					<div class="flex items-start justify-between">
						<div>
							<h3 class="font-bold text-white text-base">{branch.name}</h3>
							<p class="text-xs text-gray-400 mt-0.5">{branch.address || 'Alamat belum diatur'}</p>
						</div>
						<Badge color="green" class="font-mono text-xs">
							Radius: {branch.radius_meters}m
						</Badge>
					</div>

					<div class="bg-gray-950 rounded-xl p-3 text-xs space-y-1 font-mono text-gray-400 border border-gray-800">
						<div>Latitude: <span class="text-gray-200">{branch.latitude}</span></div>
						<div>Longitude: <span class="text-gray-200">{branch.longitude}</span></div>
					</div>

					<div class="text-[11px] text-gray-500 flex items-center justify-between pt-1 border-t border-gray-800/80">
						<span>Validasi: PostGIS Spatial SRID 4326</span>
						<span class="text-emerald-400 font-medium">✓ Geofence Aktif</span>
					</div>
				</Card>
			{/each}
		</div>
	{/if}
</div>

<!-- Modal Tambah Lokasi -->
<Modal bind:open={showModal} size="md" autoclose={false} class="bg-gray-900 text-white">
	<h3 class="text-base font-bold text-white mb-4">Tambah Lokasi Kantor / Geofence</h3>
	<form onsubmit={handleCreateBranch} class="space-y-3.5 text-xs">
		<div>
			<Label class="mb-1 text-xs text-gray-300">Nama Lokasi / Cabang</Label>
			<Input
				type="text"
				bind:value={branchName}
				required
				placeholder="cth: Kantor Cabang Surabaya / Proyek MRT"
				class="bg-gray-800 text-xs"
			/>
		</div>

		<div>
			<Label class="mb-1 text-xs text-gray-300">Alamat Lengkap</Label>
			<Textarea
				bind:value={address}
				rows="2"
				placeholder="cth: Jl. Pemuda No. 45, Genteng, Surabaya"
				class="bg-gray-800 text-xs"
			/>
		</div>

		<div class="grid grid-cols-2 gap-2">
			<div>
				<Label class="mb-1 text-xs text-gray-300">Latitude</Label>
				<Input
					type="number"
					step="any"
					bind:value={latitude}
					required
					class="bg-gray-800 text-xs"
				/>
			</div>
			<div>
				<Label class="mb-1 text-xs text-gray-300">Longitude</Label>
				<Input
					type="number"
					step="any"
					bind:value={longitude}
					required
					class="bg-gray-800 text-xs"
				/>
			</div>
		</div>

		<Button
			type="button"
			onclick={detectCurrentGPS}
			color="light"
			size="xs"
			class="w-full text-[11px]"
		>
			📍 Gunakan Titik Koordinat GPS Saat Ini
		</Button>

		<div>
			<div class="flex items-center justify-between mb-1">
				<Label class="text-xs text-gray-300">Radius Toleransi Presensi:</Label>
				<span class="font-mono text-emerald-400 font-bold">{radiusMeters} meter</span>
			</div>
			<Range
				min="10"
				max="500"
				step="5"
				bind:value={radiusMeters}
				class="accent-emerald-500 cursor-pointer"
			/>
			<p class="text-[10px] text-gray-500 mt-1">Pegawai hanya dapat melakukan presensi jika berada di dalam batas radius ini.</p>
		</div>

		<div class="flex items-center justify-end gap-2 pt-2 border-t border-gray-800">
			<Button color="light" size="sm" onclick={() => showModal = false}>Batal</Button>
			<Button type="submit" color="green" size="sm" class="font-bold text-slate-950">
				Simpan Lokasi Kantor
			</Button>
		</div>
	</form>
</Modal>
