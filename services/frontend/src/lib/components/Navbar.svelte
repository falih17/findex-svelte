<script lang="ts">
	import { Navbar, NavBrand, NavHamburger, NavUl, NavLi, Button, Badge } from 'flowbite-svelte';
	import { auth } from '$lib/stores/auth.svelte';
</script>

<Navbar fluid class="bg-gray-900 border-b border-gray-800 text-white sticky top-0 z-50 px-4 py-3">
	<NavBrand href="/">
		<div class="flex items-center gap-2">
			<span class="flex items-center justify-center w-8 h-8 rounded-lg bg-emerald-500 text-slate-950 font-black text-base shadow-md shadow-emerald-500/20">
				F
			</span>
			<span class="self-center whitespace-nowrap text-xl font-bold dark:text-white">
				Findex <Badge color="green" class="ml-1 text-[10px] font-mono">SAAS</Badge>
			</span>
		</div>
	</NavBrand>

	<div class="flex items-center md:order-2 gap-2">
		{#if auth.user}
			<div class="text-right hidden sm:block mr-2">
				<p class="text-xs font-semibold text-white leading-tight">{auth.user.name}</p>
				<Badge color={auth.user.role === 'OWNER' ? 'purple' : auth.user.role === 'ADMIN_COMPANY' ? 'blue' : 'green'} class="text-[9px] uppercase font-mono mt-0.5">
					{auth.user.role}
				</Badge>
			</div>
			<Button size="xs" color="light" class="text-xs" onclick={() => auth.logout()}>
				Keluar
			</Button>
		{:else}
			<Button size="xs" color="green" href="/login" class="text-xs font-semibold">
				Masuk
			</Button>
		{/if}
		<NavHamburger />
	</div>

	<NavUl class="text-sm font-medium">
		<NavLi href="/">Presensi (Clock-in)</NavLi>
		<NavLi href="/history">Riwayat</NavLi>

		{#if auth.user && (auth.user.role === 'OWNER' || auth.user.role === 'ADMIN_COMPANY')}
			<NavLi href="/dashboard">Dashboard HR</NavLi>
			<NavLi href="/employees">Employee</NavLi>
			<NavLi href="/shifts">Jadwal Shift</NavLi>
			<NavLi href="/branches">Lokasi Absen</NavLi>
			<NavLi href="/reports">Laporan</NavLi>
		{/if}

		<NavLi href="/download">
			<span class="text-emerald-400 font-semibold flex items-center gap-1">
				📱 Unduh APK
			</span>
		</NavLi>
	</NavUl>
</Navbar>
