<script>
	import { onMount } from 'svelte';
	import { Card, Button, Badge, Alert, Spinner } from 'flowbite-svelte';
	import { getCurrentCoordinates } from '$lib/utils/native';
	import { auth } from '$lib/stores/auth.svelte';

	let videoElement;
	let canvasElement;
	let stream = null;

	let loading = false;
	let cameraActive = false;
	let locationStatus = 'Mendeteksi lokasi...';
	let coords = { latitude: null, longitude: null, accuracy: null };
	let attendanceResult = null;
	let errorMessage = '';
	let successMessage = '';

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
			errorMessage = 'Akses kamera tidak tersedia atau belum diizinkan.';
			cameraActive = false;
		}
	}

	function stopCamera() {
		if (stream) {
			stream.getTracks().forEach((track) => track.stop());
		}
	}

	async function captureSnapshot() {
		const ctx = canvasElement.getContext('2d');
		canvasElement.width = videoElement.videoWidth || 640;
		canvasElement.height = videoElement.videoHeight || 640;
		ctx.drawImage(videoElement, 0, 0, canvasElement.width, canvasElement.height);
		return new Promise((resolve) => canvasElement.toBlob(resolve, 'image/jpeg', 0.85));
	}

	async function handleCheckIn() {
		if (!coords.latitude || !coords.longitude) {
			errorMessage = 'Lokasi GPS belum terkunci. Tunggu sebentar atau aktifkan GPS.';
			return;
		}
		if (!videoElement || !cameraActive) {
			errorMessage = 'Kamera belum aktif.';
			return;
		}

		loading = true;
		errorMessage = '';
		attendanceResult = null;

		try {
			const blob = await captureSnapshot();
			const empID = auth.user?.id || '66666666-6666-6666-6666-666666666666';

			const formData = new FormData();
			formData.append('photo', blob, 'clockin.jpg');
			formData.append('latitude', coords.latitude.toString());
			formData.append('longitude', coords.longitude.toString());
			formData.append('employee_id', empID);
			formData.append('branch_id', '22222222-2222-2222-2222-222222222222');

			const res = await fetch(`${BACKEND_URL}/attendance/check-in`, {
				method: 'POST',
				body: formData
			});

			const data = await res.json();
			if (!res.ok) {
				throw new Error(data.error || 'Presensi gagal diproses.');
			}

			attendanceResult = data.data;
			successMessage = data.message || 'Presensi masuk berhasil dicatat!';
		} catch (err) {
			errorMessage = err.message;
		} finally {
			loading = false;
		}
	}

	async function handleCheckOut() {
		loading = true;
		errorMessage = '';
		attendanceResult = null;

		try {
			const empID = auth.user?.id || '66666666-6666-6666-6666-666666666666';
			const formData = new FormData();
			formData.append('employee_id', empID);
			formData.append('latitude', (coords.latitude || 0).toString());
			formData.append('longitude', (coords.longitude || 0).toString());

			const res = await fetch(`${BACKEND_URL}/attendance/check-out`, {
				method: 'POST',
				body: formData
			});

			const data = await res.json();
			if (!res.ok) {
				throw new Error(data.error || 'Presensi pulang gagal.');
			}

			successMessage = `Presensi pulang berhasil dicatat (Status: ${data.data.check_out_status})`;
		} catch (err) {
			errorMessage = err.message;
		} finally {
			loading = false;
		}
	}

	async function handleSelfEnrollFace() {
		if (!videoElement || !cameraActive) return;
		loading = true;
		errorMessage = '';

		try {
			const blob = await captureSnapshot();
			const empID = auth.user?.id || '66666666-6666-6666-6666-666666666666';

			const formData = new FormData();
			formData.append('photo', blob, 'master.jpg');
			formData.append('employee_id', empID);

			const res = await fetch(`${BACKEND_URL}/attendance/self-enroll-face`, {
				method: 'POST',
				body: formData
			});

			const data = await res.json();
			if (!res.ok) throw new Error(data.error || 'Gagal mendaftarkan master wajah');

			successMessage = 'Wajah Anda berhasil didaftarkan sebagai basis data biometrik master!';
		} catch (err) {
			errorMessage = err.message;
		} finally {
			loading = false;
		}
	}
</script>

<div class="max-w-lg mx-auto space-y-4">
	<!-- Header -->
	<div class="text-center space-y-1">
		<div class="inline-flex items-center gap-1.5">
			<Badge color="green" class="font-mono text-xs">FINDEX BIOMETRIC</Badge>
			<Badge color="dark" class="font-mono text-xs">GEOFENCE POSTGIS</Badge>
		</div>
		<h1 class="text-2xl font-bold tracking-tight text-white">Presensi Pegawai</h1>
		<p class="text-xs text-gray-400">Verifikasi Wajah Biometrik & Validasi Geofence Kantor</p>
	</div>

	<!-- Camera Frame Card -->
	<div class="relative bg-gray-900 rounded-2xl overflow-hidden border border-gray-800 shadow-2xl aspect-square flex items-center justify-center">
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

		<!-- Camera Status Indicator Badge -->
		<div class="absolute top-3 left-3">
			<Badge color={cameraActive ? 'green' : 'red'} class="text-[11px] font-medium backdrop-blur-md">
				<span class={`w-2 h-2 rounded-full mr-1.5 ${cameraActive ? 'bg-green-400' : 'bg-red-400'}`}></span>
				{cameraActive ? 'Kamera Aktif' : 'Kamera Mati'}
			</Badge>
		</div>
	</div>

	<!-- Geolocation & Branch Info Card -->
	<Card class="bg-gray-900 border-gray-800 p-3.5 space-y-2 text-xs">
		<div class="flex items-center justify-between">
			<span class="text-gray-400 font-medium">Lokasi Kantor:</span>
			<span class="text-white font-semibold">Headquarters Jakarta</span>
		</div>
		<div class="flex items-center justify-between border-t border-gray-800 pt-2">
			<span class="text-gray-400 font-medium">Status GPS:</span>
			<span class="text-emerald-400 font-mono text-[11px]">{locationStatus}</span>
		</div>
		{#if coords.latitude}
			<div class="text-[11px] text-gray-500 font-mono">
				Lat: {coords.latitude.toFixed(5)}, Lng: {coords.longitude.toFixed(5)}
			</div>
		{/if}
	</Card>

	<!-- Error Alert -->
	{#if errorMessage}
		<Alert color="red" class="text-xs" dismissable>
			<span class="font-bold">Presensi Gagal:</span> {errorMessage}
		</Alert>
	{/if}

	<!-- Success Alert -->
	{#if successMessage}
		<Alert color="green" class="text-xs" dismissable>
			<span class="font-bold">Berhasil:</span> {successMessage}
		</Alert>
	{/if}

	<!-- Attendance Result Summary -->
	{#if attendanceResult}
		<Card class="bg-gray-900/90 border-emerald-500/30 p-4 space-y-2 text-xs">
			<div class="flex items-center justify-between border-b border-gray-800 pb-2">
				<span class="font-bold text-emerald-400 text-sm">Hasil Presensi Masuk</span>
				<Badge color={attendanceResult.status === 'ON_TIME' ? 'green' : 'yellow'} class="font-mono text-[10px]">
					{attendanceResult.status}
				</Badge>
			</div>
			<div class="grid grid-cols-2 gap-2 text-gray-300 text-[11px]">
				<div>Waktu: <span class="font-semibold text-white">{attendanceResult.check_in_time}</span></div>
				<div>Jarak Kantor: <span class="font-semibold text-white">{Math.round(attendanceResult.distance_meters || 0)}m</span></div>
				<div>Kecocokan Wajah: <span class="font-semibold text-white">{Math.round((attendanceResult.similarity_score || 0) * 100)}%</span></div>
				<div>Skor Liveness: <span class="font-semibold text-white">{Math.round((attendanceResult.liveness_score || 0) * 100)}%</span></div>
			</div>
		</Card>
	{/if}

	<!-- Action Buttons -->
	<div class="space-y-2">
		<Button
			onclick={handleCheckIn}
			disabled={loading || !cameraActive}
			color="green"
			class="w-full py-3.5 font-bold text-slate-950 text-sm shadow-lg shadow-emerald-500/20"
		>
			{#if loading}
				<Spinner class="mr-2" size="4" />
				<span>Memverifikasi Wajah & GPS...</span>
			{:else}
				<span>📸 Ambil Foto & Presensi Masuk</span>
			{/if}
		</Button>

		<div class="grid grid-cols-2 gap-2">
			<Button
				onclick={handleCheckOut}
				disabled={loading}
				color="light"
				class="py-2.5 text-xs font-semibold text-gray-300 hover:text-white"
			>
				🚪 Presensi Pulang
			</Button>

			<Button
				onclick={handleSelfEnrollFace}
				disabled={loading || !cameraActive}
				color="purple"
				class="py-2.5 text-xs font-semibold"
			>
				👤 Daftarkan Master Wajah
			</Button>
		</div>
	</div>
</div>
