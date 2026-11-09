<script>
	import { onMount } from 'svelte';
	import { getCurrentCoordinates, isNativePlatform } from '$lib/utils/native';

	let videoElement;
	let canvasElement;
	let stream = null;

	let loading = false;
	let cameraActive = false;
	let locationStatus = 'Mendeteksi lokasi...';
	let coords = { latitude: null, longitude: null, accuracy: null };
	let checkInResult = null;
	let errorMessage = '';

	const BACKEND_URL = 'http://localhost:8080/api/v1';

	onMount(() => {
		initLocation();
		startCamera();

		return () => {
			stopCamera();
		};
	});

	async function initLocation() {
		locationStatus = 'Mencari koordinat GPS...';
		try {
			const pos = await getCurrentCoordinates();
			coords = pos;
			locationStatus = `GPS Terkunci (Akurasi: ±${coords.accuracy}m)`;
		} catch (err) {
			locationStatus = `Gagal mendapatkan GPS: ${err.message}`;
		}
	}

	async function startCamera() {
		try {
			errorMessage = '';
			stream = await navigator.mediaDevices.getUserMedia({
				video: {
					facingMode: 'user',
					width: { ideal: 640 },
					height: { ideal: 640 }
				},
				audio: false
			});
			if (videoElement) {
				videoElement.srcObject = stream;
				cameraActive = true;
			}
		} catch (err) {
			console.error('Camera error:', err);
			errorMessage = 'Akses kamera ditolak atau tidak tersedia. Izinkan akses kamera untuk presensi.';
			cameraActive = false;
		}
	}

	function stopCamera() {
		if (stream) {
			stream.getTracks().forEach((track) => track.stop());
		}
	}

	async function handleCheckIn() {
		if (!coords.latitude || !coords.longitude) {
			errorMessage = 'Lokasi GPS belum didapatkan. Tunggu beberapa saat atau refresh lokasi.';
			return;
		}

		if (!videoElement || !cameraActive) {
			errorMessage = 'Kamera belum aktif.';
			return;
		}

		loading = true;
		errorMessage = '';
		checkInResult = null;

		try {
			// 1. Capture snapshot from video stream onto canvas
			const ctx = canvasElement.getContext('2d');
			canvasElement.width = videoElement.videoWidth || 640;
			canvasElement.height = videoElement.videoHeight || 640;
			ctx.drawImage(videoElement, 0, 0, canvasElement.width, canvasElement.height);

			// 2. Convert to Blob
			const blob = await new Promise((resolve) => canvasElement.toBlob(resolve, 'image/jpeg', 0.85));

			// 3. Prepare FormData
			const formData = new FormData();
			formData.append('photo', blob, 'clockin.jpg');
			formData.append('latitude', coords.latitude.toString());
			formData.append('longitude', coords.longitude.toString());
			formData.append('employee_id', '55555555-5555-5555-5555-555555555555');
			formData.append('branch_id', '22222222-2222-2222-2222-222222222222');

			// 4. Send to Backend
			const res = await fetch(`${BACKEND_URL}/attendance/check-in`, {
				method: 'POST',
				body: formData
			});

			const data = await res.json();
			if (!res.ok) {
				throw new Error(data.error || 'Presensi gagal diproses.');
			}

			checkInResult = data;
		} catch (err) {
			errorMessage = err.message;
		} finally {
			loading = false;
		}
	}
</script>

<div class="w-full space-y-4">
	<!-- Header -->
	<header class="text-center space-y-1">
		<div class="inline-flex items-center gap-2 px-3 py-1 bg-emerald-500/10 text-emerald-400 rounded-full text-xs font-medium border border-emerald-500/20">
			<span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
			Findex Multi-Tenant SaaS
		</div>
		<h1 class="text-2xl font-bold tracking-tight text-white">Presensi Masuk</h1>
		<p class="text-xs text-slate-400">Verifikasi Wajah Biometrik & Geofence Kantor</p>
	</header>

	<!-- Camera Frame Card -->
	<div class="relative bg-slate-800 rounded-2xl overflow-hidden border border-slate-700/80 shadow-2xl aspect-square flex items-center justify-center">
		<!-- Live Video Stream -->
		<video
			bind:this={videoElement}
			autoplay
			playsinline
			muted
			class="w-full h-full object-cover -scale-x-100"
		></video>

		<!-- Hidden Canvas for Snapshot Capture -->
		<canvas bind:this={canvasElement} class="hidden"></canvas>

		<!-- Face Oval Silhouette Overlay Guide -->
		<div class="absolute inset-0 pointer-events-none flex items-center justify-center">
			<div class="w-48 h-64 border-2 border-dashed border-emerald-400/60 rounded-[50%] animate-pulse"></div>
		</div>

		<!-- Camera Status Indicator -->
		<div class="absolute top-3 left-3 bg-black/60 backdrop-blur-md px-2.5 py-1 rounded-full text-[11px] font-medium text-slate-300 flex items-center gap-1.5">
			<span class={`w-2 h-2 rounded-full ${cameraActive ? 'bg-green-500' : 'bg-red-500'}`}></span>
			{cameraActive ? 'Kamera Aktif' : 'Kamera Mati'}
		</div>
	</div>

	<!-- Geolocation & Branch Info -->
	<div class="bg-slate-800/80 backdrop-blur-md rounded-xl p-3.5 border border-slate-700/60 space-y-2 text-xs">
		<div class="flex items-center justify-between">
			<span class="text-slate-400 font-medium">Lokasi Kantor Tujuan:</span>
			<span class="text-white font-semibold">Headquarters Jakarta</span>
		</div>
		<div class="flex items-center justify-between border-t border-slate-700/40 pt-2">
			<span class="text-slate-400 font-medium">Status GPS Anda:</span>
			<span class="text-emerald-400 font-mono text-[11px]">{locationStatus}</span>
		</div>
		{#if coords.latitude}
			<div class="text-[11px] text-slate-500 font-mono">
				Lat: {coords.latitude.toFixed(5)}, Lng: {coords.longitude.toFixed(5)}
			</div>
		{/if}
	</div>

	<!-- Error Alert -->
	{#if errorMessage}
		<div class="p-3 bg-red-500/10 border border-red-500/20 text-red-400 rounded-xl text-xs flex items-start gap-2">
			<span class="text-base">⚠️</span>
			<div>
				<p class="font-semibold">Presensi Gagal</p>
				<p>{errorMessage}</p>
			</div>
		</div>
	{/if}

	<!-- Success Alert -->
	{#if checkInResult}
		<div class="p-4 bg-emerald-500/10 border border-emerald-500/30 text-emerald-300 rounded-xl text-xs space-y-1.5">
			<div class="flex items-center gap-2 text-sm font-bold text-emerald-400">
				<span>✅</span> Presensi Berhasil Dicatat!
			</div>
			<div class="grid grid-cols-2 gap-1 text-[11px] pt-1 border-t border-emerald-500/20">
				<div>Waktu: <span class="font-semibold text-white">{checkInResult.time}</span></div>
				<div>Status: <span class="font-semibold text-white">{checkInResult.status}</span></div>
				<div>Jarak Kantor: <span class="font-semibold text-white">{Math.round(checkInResult.distance_meters || 0)} meter</span></div>
				<div>Liveness: <span class="font-semibold text-white">{Math.round((checkInResult.liveness_score || 0) * 100)}%</span></div>
			</div>
		</div>
	{/if}

	<!-- Action Button -->
	<button
		onclick={handleCheckIn}
		disabled={loading || !cameraActive}
		class="w-full py-3.5 px-4 bg-emerald-500 hover:bg-emerald-400 active:scale-[0.98] transition-all disabled:opacity-50 disabled:pointer-events-none rounded-xl font-semibold text-slate-950 text-sm shadow-lg shadow-emerald-500/25 flex items-center justify-center gap-2"
	>
		{#if loading}
			<svg class="animate-spin h-5 w-5 text-slate-950" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
			<span>Memverifikasi Wajah & Lokasi...</span>
		{:else}
			<span>📸 Ambil Foto & Presensi Sekarang</span>
		{/if}
	</button>
</div>
