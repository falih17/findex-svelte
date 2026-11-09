<script lang="ts">
	import { Card, Badge, Button, Alert } from 'flowbite-svelte';

	let apkVersion = '1.0.0';
	let releaseDate = '2026-11-09';
	let fileSize = '18.4 MB';
</script>

<div class="max-w-xl mx-auto space-y-6 py-6">
	<!-- App Icon & Badge -->
	<div class="text-center space-y-3">
		<div class="inline-flex items-center justify-center w-20 h-20 rounded-2xl bg-gradient-to-tr from-emerald-600 to-teal-400 text-slate-950 font-black text-3xl shadow-xl shadow-emerald-500/25 mx-auto border-2 border-emerald-400/30">
			F
		</div>
		<div>
			<h1 class="text-2xl font-bold text-white tracking-tight">Unduh Aplikasi Findex</h1>
			<p class="text-xs text-gray-400 mt-1">Khusus Karyawan Perusahaan untuk Presensi Wajah & GPS Native</p>
		</div>
		<div class="inline-flex items-center gap-2">
			<Badge color="green" class="font-mono text-xs">Android APK v{apkVersion}</Badge>
			<Badge color="dark" class="font-mono text-xs">{fileSize}</Badge>
		</div>
	</div>

	<!-- Download Card -->
	<Card padding="lg" class="bg-gray-900 border-gray-800 space-y-4 max-w-none shadow-xl">
		<h2 class="text-sm font-bold text-white flex items-center gap-2">
			<span>📲</span> Keunggulan Aplikasi Native Android:
		</h2>
		<ul class="text-xs text-gray-400 space-y-2.5 list-disc list-inside">
			<li><strong class="text-gray-200">Kamera Depan Ultra-Cepat:</strong> Verifikasi biometrik wajah langsung tanpa lag dengan anti-spoofing AI.</li>
			<li><strong class="text-gray-200">Penguncian GPS Presisi:</strong> Menggunakan hardware GPS ponsel dengan validasi polygon radius kantor.</li>
			<li><strong class="text-gray-200">Proteksi Anti-Fake GPS:</strong> Deteksi mock location dan manipulasi provider posisi.</li>
			<li><strong class="text-gray-200">Notifikasi Jadwal:</strong> Pengingat otomatis menjelang batas check-in dan waktu awal checkout.</li>
		</ul>

		<!-- Download Action Button -->
		<div class="pt-2">
			<Button
				size="lg"
				color="green"
				class="w-full font-bold text-slate-950 flex items-center justify-center gap-2 shadow-lg shadow-emerald-500/20"
				href="/findex-attendance.apk"
				download="findex-attendance.apk"
			>
				<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
				</svg>
				<span>Unduh APK Android (Langsung)</span>
			</Button>
		</div>
	</Card>

	<!-- Instructions -->
	<Card padding="md" class="bg-gray-900/60 border-gray-800 text-xs space-y-2 text-gray-400 max-w-none">
		<p class="font-semibold text-gray-200 flex items-center gap-1.5">
			<span>💡</span> Panduan Instalasi:
		</p>
		<ol class="list-decimal list-inside space-y-1.5 text-[11px] leading-relaxed text-gray-400">
			<li>Klik tombol hijau di atas untuk mengunduh file <code class="text-emerald-400 font-mono">findex-attendance.apk</code>.</li>
			<li>Buka file hasil unduhan pada file manager atau browser ponsel Android Anda.</li>
			<li>Jika muncul peringatan keamanan, pilih <em>"Izinkan instalasi dari sumber ini"</em> (Allow from this source).</li>
			<li>Tekan <strong>Install</strong> dan buka aplikasi. Berikan izin <strong>Kamera</strong> & <strong>Lokasi</strong> saat diminta.</li>
			<li>Masuk menggunakan akun email dan kata sandi yang telah diberikan oleh HRD perusahaan.</li>
		</ol>
	</Card>
</div>