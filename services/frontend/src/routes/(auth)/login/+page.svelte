<script lang="ts">
	import { Card, Label, Input, Button, Alert, Spinner } from 'flowbite-svelte';
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

<div class="flex justify-center items-center py-8">
	<Card class="w-full max-w-md bg-gray-900 border-gray-800 shadow-2xl p-6 sm:p-8">
		<div class="text-center space-y-2 mb-6">
			<div class="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-emerald-500 text-slate-950 font-black text-2xl shadow-lg shadow-emerald-500/20">
				F
			</div>
			<h1 class="text-2xl font-bold text-white tracking-tight">Masuk ke Findex</h1>
			<p class="text-xs text-gray-400">SaaS Presensi Wajah & Geofence Multi-Tenant</p>
		</div>

		{#if errorMsg}
			<Alert color="red" class="mb-4 text-xs" dismissable>
				<span class="font-medium">Gagal:</span> {errorMsg}
			</Alert>
		{/if}

		<form onsubmit={handleLogin} class="space-y-4">
			<div>
				<Label for="email" class="mb-2 text-xs text-gray-300">Email Akun</Label>
				<Input
					id="email"
					type="email"
					bind:value={email}
					required
					placeholder="nama@perusahaan.com"
					class="bg-gray-800 border-gray-700 text-white text-xs sm:text-sm"
				/>
			</div>

			<div>
				<Label for="password" class="mb-2 text-xs text-gray-300">Kata Sandi</Label>
				<Input
					id="password"
					type="password"
					bind:value={password}
					required
					class="bg-gray-800 border-gray-700 text-white text-xs sm:text-sm"
				/>
			</div>

			<!-- Quick Role Switcher for Testing -->
			<div class="p-3 bg-gray-800/60 rounded-xl border border-gray-700/60 text-xs text-gray-400 space-y-2">
				<p class="font-semibold text-gray-300">Pilih Role Akun Demo:</p>
				<div class="flex flex-wrap gap-1.5">
					<button
						type="button"
						onclick={() => setRoleCredential('owner@findex.local')}
						class="px-2 py-1 bg-purple-500/10 hover:bg-purple-500/20 text-purple-400 rounded border border-purple-500/30 text-[11px]"
					>
						👑 Owner
					</button>
					<button
						type="button"
						onclick={() => setRoleCredential('admin@findex.local')}
						class="px-2 py-1 bg-blue-500/10 hover:bg-blue-500/20 text-blue-400 rounded border border-blue-500/30 text-[11px]"
					>
						🏢 Admin Company
					</button>
					<button
						type="button"
						onclick={() => setRoleCredential('employee@findex.local')}
						class="px-2 py-1 bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400 rounded border border-emerald-500/30 text-[11px]"
					>
						👤 Employee
					</button>
				</div>
				<p class="text-[10px] text-gray-500">Kata Sandi: <code class="text-emerald-400 font-mono">password123</code></p>
			</div>

			<Button type="submit" color="green" class="w-full font-bold text-slate-950" disabled={loading}>
				{#if loading}
					<Spinner class="mr-3" size="4" />
					<span>Memverifikasi...</span>
				{:else}
					<span>Masuk Sekarang</span>
				{/if}
			</Button>
		</form>
	</Card>
</div>
