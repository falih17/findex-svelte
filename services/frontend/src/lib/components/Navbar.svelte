<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
</script>

<nav class="w-full bg-slate-900/90 backdrop-blur-md border-b border-slate-800 sticky top-0 z-50 px-4 py-3">
	<div class="max-w-6xl mx-auto flex items-center justify-between">
		<!-- Brand Logo -->
		<a href="/" class="flex items-center gap-2 font-bold text-lg text-white tracking-tight">
			<span class="flex items-center justify-center w-8 h-8 rounded-lg bg-emerald-500 text-slate-950 font-black text-base shadow-md shadow-emerald-500/20">
				F
			</span>
			<span>Findex <span class="text-emerald-400 font-normal text-xs uppercase tracking-widest ml-1">SaaS</span></span>
		</a>

		<!-- Navigation Links -->
		<div class="flex items-center gap-1 sm:gap-3 text-xs sm:text-sm font-medium">
			<a
				href="/"
				class="px-3 py-1.5 rounded-lg text-slate-300 hover:text-white hover:bg-slate-800 transition-colors"
			>
				Clock-in
			</a>
			<a
				href="/history"
				class="px-3 py-1.5 rounded-lg text-slate-300 hover:text-white hover:bg-slate-800 transition-colors"
			>
				Riwayat
			</a>
			<a
				href="/dashboard"
				class="px-3 py-1.5 rounded-lg text-emerald-400 hover:bg-emerald-500/10 transition-colors border border-emerald-500/20"
			>
				HR Portal
			</a>
			<a
				href="/branches"
				class="hidden sm:inline-block px-3 py-1.5 rounded-lg text-slate-300 hover:text-white hover:bg-slate-800 transition-colors"
			>
				Lokasi Kantor
			</a>
			<a
				href="/download"
				class="px-2.5 py-1 rounded-lg bg-emerald-500/10 text-emerald-400 hover:bg-emerald-500/20 transition-colors text-xs font-semibold border border-emerald-500/20"
			>
				📱 Unduh APK
			</a>
		</div>

		<!-- Auth Status / Demo Badge -->
		<div class="flex items-center gap-2">
			{#if auth.user}
				<div class="text-right hidden sm:block">
					<p class="text-xs font-semibold text-white leading-tight">{auth.user.name}</p>
					<p class="text-[10px] text-emerald-400 font-mono">{auth.user.role}</p>
				</div>
				<button
					onclick={() => auth.logout()}
					class="text-xs text-slate-400 hover:text-red-400 px-2 py-1 rounded bg-slate-800"
				>
					Logout
				</button>
			{:else}
				<a
					href="/login"
					class="text-xs bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-semibold px-3 py-1.5 rounded-lg transition-all"
				>
					Masuk
				</a>
			{/if}
		</div>
	</div>
</nav>
