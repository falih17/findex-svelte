<script lang="ts">
	let apkVersion = '1.0.0';
	let releaseDate = '2026-11-09';
	let fileSize = '18.4 MB';
</script>

<div class="max-w-md mx-auto space-y-6 text-center py-6">
	<!-- App Icon & Badge -->
	<div class="space-y-3">
		<div class="inline-flex items-center justify-center w-20 h-20 rounded-2xl bg-gradient-to-tr from-emerald-600 to-teal-400 text-slate-950 font-black text-3xl shadow-xl shadow-emerald-500/25 mx-auto border-2 border-emerald-400/30">
			F
		</div>
		<div>
			<h1 class="text-2xl font-bold text-white tracking-tight">Unduh Aplikasi Findex</h1>
			<p class="text-xs text-slate-400 mt-1">Khusus Karyawan Perusahaan untuk Presensi Wajah & GPS</p>
		</div>
		<div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-slate-800 text-slate-300 text-xs font-mono border border-slate-700">
			<span>Android APK v{apkVersion}</span> • <span>{fileSize}</span>
		</div>
	</div>

	<!-- Download Card -->
	<div class="bg-slate-900 border border-slate-800 rounded-2xl p-6 space-y-4 shadow-2xl text-left">
		<h2 class="text-sm font-bold text-white flex items-center gap-2">
			<span>📲</span> Keunggulan Aplikasi Native Android:
		</h2>
		<ul class="text-xs text-slate-400 space-y-2 list-disc list-inside">
			<li>Akses kamera depan ultra-cepat & jernih untuk verifikasi wajah biometrik.</li>
			<li>Penguncian koordinat GPS real-time dengan akurasi tinggi (High-Accuracy).</li>
			<li>Proteksi anti-fake GPS dan anti-spoofing terintegrasi.</li>
			<li>Notifikasi pengingat jam masuk dan jam pulang kerja otomatis.</li>
		</ul>

		<!-- Download Action Button -->
		<a
			href="/findex-attendance.apk"
			download="findex-attendance.apk"
			class="w-full py-3.5 px-4 bg-emerald-500 hover:bg-emerald-400 active:scale-[0.98] transition-all rounded-xl font-bold text-slate-950 text-sm shadow-lg shadow-emerald-500/25 flex items-center justify-center gap-2 mt-4"
		>
			<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
			</svg>
			<span>Unduh APK Android (Langsung)</span>
		</a>
	</div>

	<!-- Instructions -->
	<div class="bg-slate-900/60 border border-slate-800/80 rounded-xl p-4 text-left text-xs space-y-2 text-slate-400">
		<p class="font-semibold text-slate-300">💡 Panduan Instalasi:</p>
		<ol class="list-decimal list-inside space-y-1 text-[11px] leading-relaxed">
			<li>Klik tombol hijau di atas untuk mengunduh file <code class="text-emerald-400 font-mono">findex-attendance.apk</code>.</li>
			<li>Buka file hasil unduhan pada ponsel Android Anda.</li>
			<li>Jika muncul peringatan keamanan browser, pilih <em>"Izinkan instalasi dari sumber ini"</em> (Allow from this source).</li>
			<li>Tekan <strong>Install</strong> dan buka aplikasi. Berikan izin Akses Kamera & Lokasi saat pertama kali dibuka.</li>
			<li>Masuk menggunakan email dan kata sandi yang telah didaftarkan oleh HRD perusahaan Anda.</li>
		</ol>
	</div>
</div>
