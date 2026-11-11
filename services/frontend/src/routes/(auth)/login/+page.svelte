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

			if (json.data.employee.role === 'EMPLOYEE') {
				goto('/');
			} else {
				goto('/dashboard');
			}
		} catch (err: any) {
			errorMsg = err.message;
		} finally {
			loading = false;
		}
	}

	function setRoleCredential(targetEmail: string) {
		email = targetEmail;
		password = 'password123';
	}
</script>

<div class="flex justify-center items-center py-10">
	<div class="w-full max-w-md bg-white border border-slate-200/90 rounded-2xl shadow-xl p-6 sm:p-8 space-y-6">
		<div class="text-center space-y-2">
			<!-- AbsenPro / Findex Icon -->
			<div class="w-12 h-12 rounded-xl bg-indigo-600 flex items-center justify-center text-white mx-auto shadow-md shadow-indigo-600/30">
				<svg class="w-7 h-7" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path stroke-linecap="round" stroke-linejoin="round" d="M3 8V5a2 2 0 012-2h3m8 0h3a2 2 0 012 2v3m0 8v3a2 2 0 01-2 2h-3m-8 0H5a2 2 0 01-2-2v-3" />
					<circle cx="9" cy="9" r="1.5" fill="currentColor" />
					<circle cx="15" cy="9" r="1.5" fill="currentColor" />
					<path stroke-linecap="round" stroke-linejoin="round" d="M8.5 14.5c1 .9 2 1.4 3.5 1.4s2.5-.5 3.5-1.4" />
				</svg>
			</div>
			<div>
				<h1 class="text-2xl font-bold tracking-tight text-slate-900">
					Masuk ke <span class="text-indigo-600">AbsenPro</span>
				</h1>
				<p class="text-xs text-slate-500 mt-0.5">SaaS Presensi Biometrik Wajah & Geofence</p>
			</div>
		</div>

		{#if errorMsg}
			<div class="p-3 bg-rose-50 border border-rose-200 text-rose-700 rounded-xl text-xs flex items-center gap-2">
				<span>⚠️</span>
				<span>{errorMsg}</span>
			</div>
		{/if}

		<form onsubmit={handleLogin} class="space-y-4 text-xs">
			<div>
				<label for="email" class="block mb-1.5 font-semibold text-slate-700">Email Akun</label>
				<input
					id="email"
					type="email"
					bind:value={email}
					required
					placeholder="nama@perusahaan.com"
					class="w-full px-3.5 py-2.5 bg-slate-50 border border-slate-200 rounded-lg text-slate-800 text-xs sm:text-sm focus:ring-2 focus:ring-indigo-500 focus:bg-white transition-all"
				/>
			</div>

			<div>
				<label for="password" class="block mb-1.5 font-semibold text-slate-700">Kata Sandi</label>
				<input
					id="password"
					type="password"
					bind:value={password}
					required
					class="w-full px-3.5 py-2.5 bg-slate-50 border border-slate-200 rounded-lg text-slate-800 text-xs sm:text-sm focus:ring-2 focus:ring-indigo-500 focus:bg-white transition-all"
				/>
			</div>

			<!-- Quick Role Switcher for Testing -->
			<div class="p-3 bg-slate-50 rounded-xl border border-slate-200/80 space-y-2">
				<p class="font-semibold text-slate-700 text-[11px]">Pilih Role Akun Demo:</p>
				<div class="flex flex-wrap gap-1.5">
					<button
						type="button"
						onclick={() => setRoleCredential('owner@findex.local')}
						class="px-2.5 py-1 bg-purple-50 hover:bg-purple-100 text-purple-700 rounded-md border border-purple-200 text-[11px] font-semibold transition-colors"
					>
						👑 Owner
					</button>
					<button
						type="button"
						onclick={() => setRoleCredential('admin@findex.local')}
						class="px-2.5 py-1 bg-indigo-50 hover:bg-indigo-100 text-indigo-700 rounded-md border border-indigo-200 text-[11px] font-semibold transition-colors"
					>
						🏢 Admin Company
					</button>
					<button
						type="button"
						onclick={() => setRoleCredential('employee@findex.local')}
						class="px-2.5 py-1 bg-emerald-50 hover:bg-emerald-100 text-emerald-700 rounded-md border border-emerald-200 text-[11px] font-semibold transition-colors"
					>
						👤 Employee
					</button>
				</div>
				<p class="text-[10px] text-slate-400">Kata Sandi: <code class="text-indigo-600 font-mono">password123</code></p>
			</div>

			<button
				type="submit"
				disabled={loading}
				class="w-full py-2.5 px-4 bg-indigo-600 hover:bg-indigo-700 active:scale-[0.99] disabled:opacity-50 text-white font-bold text-xs sm:text-sm rounded-lg shadow-md shadow-indigo-600/25 transition-all flex items-center justify-center gap-2"
			>
				{#if loading}
					<span>Memverifikasi Akun...</span>
				{:else}
					<span>Masuk Sekarang</span>
				{/if}
			</button>
		</form>
	</div>
</div>