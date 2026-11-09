<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';

	let email = $state('admin@findex.local');
	let password = $state('password123');
	let loading = $state(false);
	let errorMsg = $state('');

	const BACKEND_URL = 'http://localhost:8080/api/v1';

	async function handleLogin(e: Event) {
		e.preventDefault();
		loading = true;
		errorMsg = '';

		try {
			const res = await fetch(`${BACKEND_URL}/auth/login`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email, password })
			});

			const json = await res.json();
			if (!res.ok) {
				throw new Error(json.error || 'Login gagal.');
			}

			auth.setSession(json.data.token, json.data.employee);
			goto('/');
		} catch (err: any) {
			errorMsg = err.message;
		} finally {
			loading = false;
		}
	}
</script>

<div class="max-w-md mx-auto mt-8 bg-slate-900 border border-slate-800 rounded-2xl p-6 sm:p-8 shadow-2xl">
	<div class="text-center space-y-2 mb-6">
		<div class="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-emerald-500 text-slate-950 font-black text-2xl shadow-lg shadow-emerald-500/20">
			F
		</div>
		<h1 class="text-2xl font-bold text-white tracking-tight">Masuk ke Findex</h1>
		<p class="text-xs text-slate-400">Portal Presensi & HR Management Terintegrasi</p>
	</div>

	{#if errorMsg}
		<div class="mb-4 p-3 bg-red-500/10 border border-red-500/20 text-red-400 rounded-xl text-xs">
			{errorMsg}
		</div>
	{/if}

	<form onsubmit={handleLogin} class="space-y-4 text-sm">
		<div>
			<label for="email" class="block text-xs font-medium text-slate-300 mb-1">Email Perusahaan</label>
			<input
				id="email"
				type="email"
				bind:value={email}
				required
				class="w-full px-3.5 py-2.5 bg-slate-800/80 border border-slate-700 rounded-xl text-white focus:outline-none focus:border-emerald-500 text-xs sm:text-sm"
			/>
		</div>

		<div>
			<label for="password" class="block text-xs font-medium text-slate-300 mb-1">Kata Sandi</label>
			<input
				id="password"
				type="password"
				bind:value={password}
				required
				class="w-full px-3.5 py-2.5 bg-slate-800/80 border border-slate-700 rounded-xl text-white focus:outline-none focus:border-emerald-500 text-xs sm:text-sm"
			/>
		</div>

		<div class="p-3 bg-slate-800/50 rounded-lg border border-slate-700/50 text-[11px] text-slate-400">
			💡 <strong>Akun Demo Otomatis:</strong><br />
			Email: <code class="text-emerald-400 font-mono">admin@findex.local</code><br />
			Password: <code class="text-emerald-400 font-mono">password123</code>
		</div>

		<button
			type="submit"
			disabled={loading}
			class="w-full py-3 bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-semibold rounded-xl text-sm transition-all disabled:opacity-50"
		>
			{loading ? 'Memverifikasi...' : 'Masuk Sekarang'}
		</button>
	</form>
</div>
