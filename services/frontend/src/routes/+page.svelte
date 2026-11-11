<script>
	import { onMount } from 'svelte';
	import { getCurrentCoordinates } from '$lib/utils/native';
	import { auth } from '$lib/stores/auth.svelte';

	let videoElement;
	let canvasElement;
	let stream = null;

	let loading = false;
	let cameraActive = false;
	let locationStatus = 'Mendeteksi lokasi GPS...';
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

<div class="max-w-xl mx-auto space-y-4">
	<!-- Header -->
	<div class="text-center space-y-1">
		<div class="inline-flex items-center gap-1.5">
			<span class="px-2.5 py-0.5 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-700 font-mono text-[11px] font-bold">
				BIOMETRIC AI
			</span>
			<span class="px-2.5 py-0.5 rounded-full bg-indigo-50 border border-indigo-200 text-indigo-700 font-mono text-[11px] font-bold">
				GEOFENCE POSTGIS
			</span>
		</div>
		<h1 class="text-2xl font-bold tracking-tight text-slate-900">Presensi & Log Live</h1>
		<p class="text-xs text-slate-500">Verifikasi Wajah Biometrik & Validasi Geofence Radius Kantor</p>
	</div>

	<!-- Camera Frame Card -->
	<div class="relative bg-slate-950 rounded-2xl overflow-hidden border border-slate-300 shadow-xl aspect-square flex items-center justify-center">
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
			<div class="w-48 h-64 border-2 border-dashed border-emerald-400/70 rounded-[50%] animate-pulse"></div>
		</div>

		<!-- Camera Status Indicator Badge -->
		<div class="absolute top-3 left-3">
			<div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-slate-900/80 backdrop-blur-md text-white text-xs font-medium border border-slate-700">
				<span class={`w-2 h-2 rounded-full ${cameraActive ? 'bg-emerald-400 animate-pulse' : 'bg-rose-500'}`}></span>
				<span>{cameraActive ? 'Kamera Siap' : 'Kamera Mati'}</span>
			</div>
		</div>
	</div>

	<!-- Geolocation & Branch Info Card -->
	<div class="bg-white border border-slate-200/90 rounded-2xl p-4 space-y-2 text-xs shadow-2xs">
		<div class="flex items-center justify-between">
			<span class="text-slate-500 font-medium">Lokasi Kantor:</span>
			<span class="text-slate-900 font-bold">Headquarters Jakarta (Radius 100m)</span>
		</div>
		<div class="flex items-center justify-between border-t border-slate-100 pt-2">
			<span class="text-slate-500 font-medium">Status GPS:</span>
			<span class="text-emerald-700 font-mono font-semibold text-[11px] bg-emerald-50 px-2 py-0.5 rounded border border-emerald-200">{locationStatus}</span>
		</div>
		{#if coords.latitude}
			<div class="text-[11px] text-slate-400 font-mono">
				Koordinat: {coords.latitude.toFixed(5)}, {coords.longitude.toFixed(5)}
			</div>
		{/if}
	</div>

	<!-- Error Alert -->
	{#if errorMessage}
		<div class="p-3 bg-rose-50 border border-rose-200 text-rose-700 rounded-xl text-xs flex items-center gap-2">
			<span>⚠️</span>
			<span>{errorMessage}</span>
		</div>
	{/if}

	<!-- Success Alert -->
	{#if successMessage}
		<div class="p-3 bg-emerald-50 border border-emerald-200 text-emerald-800 rounded-xl text-xs flex items-center gap-2 font-medium">
			<span>✅</span>
			<span>{successMessage}</span>
		</div>
	{/if}

	<!-- Attendance Result Summary -->
	{#if attendanceResult}
		<div class="bg-white border border-emerald-300 rounded-2xl p-4 space-y-2 text-xs shadow-sm">
			<div class="flex items-center justify-between border-b border-slate-100 pb-2">
				<span class="font-bold text-emerald-700 text-sm">Hasil Presensi Masuk</span>
				<span class="px-2 py-0.5 rounded-full bg-emerald-100 text-emerald-800 font-mono font-bold text-[10px]">
					{attendanceResult.status}
				</span>
			</div>
			<div class="grid grid-cols-2 gap-2 text-slate-600 text-[11px]">
				<div>Waktu: <span class="font-bold text-slate-900">{attendanceResult.check_in_time}</span></div>
				<div>Jarak Kantor: <span class="font-bold text-slate-900">{Math.round(attendanceResult.distance_meters || 0)}m</span></div>
				<div>Kecocokan Wajah: <span class="font-bold text-emerald-600">{Math.round((attendanceResult.similarity_score || 0) * 100)}%</span></div>
				<div>Skor Liveness: <span class="font-bold text-emerald-600">{Math.round((attendanceResult.liveness_score || 0) * 100)}%</span></div>
			</div>
		</div>
	{/if}

	<!-- Action Buttons -->
	<div class="space-y-2.5">
		<button
			type="button"
			onclick={handleCheckIn}
			disabled={loading || !cameraActive}
			class="w-full py-3.5 bg-indigo-600 hover:bg-indigo-700 active:scale-[0.99] disabled:opacity-50 text-white font-bold text-sm rounded-xl shadow-md shadow-indigo-600/25 flex items-center justify-center gap-2 transition-all"
		>
			{#if loading}
				<span>Memverifikasi Wajah & GPS...</span>
			{:else}
				<span>📸 Ambil Foto & Presensi Masuk</span>
			{/if}
		</button>

		<div class="grid grid-cols-2 gap-2">
			<button
				type="button"
				onclick={handleCheckOut}
				disabled={loading}
				class="py-2.5 bg-white hover:bg-slate-50 border border-slate-200 text-slate-700 text-xs font-semibold rounded-xl shadow-2xs transition-colors"
			>
				🚪 Presensi Pulang
			</button>

			<button
				type="button"
				onclick={handleSelfEnrollFace}
				disabled={loading || !cameraActive}
				class="py-2.5 bg-indigo-50 hover:bg-indigo-100 border border-indigo-200 text-indigo-700 text-xs font-semibold rounded-xl transition-colors"
			>
				👤 Daftarkan Master Wajah
			</button>
		</div>
	</div>
</div>