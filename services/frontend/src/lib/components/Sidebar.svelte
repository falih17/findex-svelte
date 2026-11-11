<script lang="ts">
	import { page } from '$app/stores';
	import { auth } from '$lib/stores/auth.svelte';

	let { open = $bindable(false) } = $props();

	// Navigation items matching AbsenPro screenshot
	const navItems = [
		{
			label: 'Dashboard',
			href: '/dashboard',
			icon: 'dashboard',
			activePattern: /^\/dashboard/,
			indicator: 'active-dot'
		},
		{
			label: 'Presensi & Log Live',
			href: '/',
			icon: 'map-pin',
			activePattern: /^\/$/,
			badge: { text: 'Live', color: 'green' }
		},
		{
			label: 'Persetujuan',
			href: '/reports',
			icon: 'check-circle',
			activePattern: /^\/approvals/,
			badge: { text: '5', color: 'amber' }
		},
		{
			label: 'Karyawan & Divisi',
			href: '/employees',
			icon: 'users',
			activePattern: /^\/employees/
		},
		{
			label: 'Geofencing & Cabang',
			href: '/branches',
			icon: 'crosshair',
			activePattern: /^\/branches/
		},
		{
			label: 'Jadwal & Shift Kerja',
			href: '/shifts',
			icon: 'calendar',
			activePattern: /^\/shifts/
		},
		{
			label: 'Laporan & Payroll',
			href: '/reports',
			icon: 'file-text',
			activePattern: /^\/reports/
		},
		{
			label: 'Riwayat Presensi',
			href: '/history',
			icon: 'history',
			activePattern: /^\/history/
		},
		{
			label: 'Unduh APK Mobile',
			href: '/download',
			icon: 'smartphone',
			activePattern: /^\/download/
		}
	];

	const adminItems = [
		{
			label: 'Manajemen Tenant',
			href: '/companies',
			icon: 'building',
			ownerOnly: true,
			badge: { text: 'Owner', color: 'purple' }
		},
		{
			label: 'Pengaturan Tenant',
			href: '/companies',
			icon: 'settings'
		},
		{
			label: 'Audit Trail & Anti-Fraud',
			href: '#',
			icon: 'shield-check'
		}
	];

	function isItemActive(pattern?: RegExp) {
		if (!pattern) return false;
		return pattern.test($page.url.pathname);
	}
</script>

<!-- Backdrop on mobile -->
{#if open}
	<div 
		class="fixed inset-0 bg-slate-950/70 z-40 lg:hidden backdrop-blur-xs transition-opacity"
		onclick={() => open = false}
	></div>
{/if}

<!-- Sidebar Container -->
<aside
	class={`fixed top-0 bottom-0 left-0 z-50 w-64 bg-[#0c1322] text-slate-300 flex flex-col border-r border-slate-800/80 transition-transform duration-300 lg:translate-x-0 ${
		open ? 'translate-x-0' : '-translate-x-full'
	}`}
>
	<!-- Top Brand / Logo Header -->
	<div class="px-5 py-4 border-b border-slate-800/60 flex items-center justify-between">
		<a href="/dashboard" class="flex items-center gap-3 group">
			<!-- Face Scan / Biometric Icon -->
			<div class="w-10 h-10 rounded-xl bg-indigo-600 flex items-center justify-center text-white shadow-lg shadow-indigo-600/30">
				<svg class="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path stroke-linecap="round" stroke-linejoin="round" d="M3 8V5a2 2 0 012-2h3m8 0h3a2 2 0 012 2v3m0 8v3a2 2 0 01-2 2h-3m-8 0H5a2 2 0 01-2-2v-3" />
					<circle cx="9" cy="9" r="1.5" fill="currentColor" />
					<circle cx="15" cy="9" r="1.5" fill="currentColor" />
					<path stroke-linecap="round" stroke-linejoin="round" d="M8.5 14.5c1 .9 2 1.4 3.5 1.4s2.5-.5 3.5-1.4" />
				</svg>
			</div>
			<div>
				<div class="text-base font-bold text-white tracking-tight flex items-center">
					<span>Absen</span><span class="text-indigo-400">Pro</span>
				</div>
				<div class="text-[9px] font-semibold tracking-wider text-slate-400 uppercase">
					Cloud SaaS Platform
				</div>
			</div>
		</a>

		<!-- Mobile Close Button -->
		<button
			type="button"
			class="lg:hidden text-slate-400 hover:text-white p-1"
			onclick={() => open = false}
			aria-label="Tutup menu"
		>
			<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
			</svg>
		</button>
	</div>

	<!-- Tenant Selector -->
	<div class="px-4 py-3 border-b border-slate-800/40">
		<div class="text-[10px] font-bold tracking-wider text-slate-400 uppercase mb-1.5 px-1">
			Tenant Aktif
		</div>
		<div class="w-full bg-[#141d33] border border-slate-800 rounded-lg px-3 py-2 flex items-center justify-between text-xs text-slate-200 cursor-pointer hover:border-slate-700 transition-colors">
			<span class="font-medium truncate pr-2">PT Maju Bersama Makmur (HQ)</span>
			<svg class="w-3.5 h-3.5 text-slate-400 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l4-4 4 4m0 6l-4 4-4-4" />
			</svg>
		</div>
	</div>

	<!-- Navigation Menu List -->
	<nav class="flex-1 overflow-y-auto px-3 py-3 space-y-1 text-xs">
		{#each navItems as item}
			{@const active = isItemActive(item.activePattern)}
			<a
				href={item.href}
				onclick={() => open = false}
				class={`flex items-center justify-between px-3 py-2.5 rounded-lg font-medium transition-all group ${
					active
						? 'bg-indigo-600 text-white shadow-sm shadow-indigo-600/30'
						: 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50'
				}`}
			>
				<div class="flex items-center gap-3">
					<!-- Icons -->
					{#if item.icon === 'dashboard'}
						<svg class={`w-4 h-4 ${active ? 'text-white' : 'text-slate-400 group-hover:text-slate-200'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
						</svg>
					{:else if item.icon === 'map-pin'}
						<svg class={`w-4 h-4 ${active ? 'text-white' : 'text-slate-400 group-hover:text-slate-200'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
						</svg>
					{:else if item.icon === 'check-circle'}
						<svg class={`w-4 h-4 ${active ? 'text-white' : 'text-slate-400 group-hover:text-slate-200'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					{:else if item.icon === 'users'}
						<svg class={`w-4 h-4 ${active ? 'text-white' : 'text-slate-400 group-hover:text-slate-200'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
						</svg>
					{:else if item.icon === 'crosshair'}
						<svg class={`w-4 h-4 ${active ? 'text-white' : 'text-slate-400 group-hover:text-slate-200'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<circle cx="12" cy="12" r="8" stroke-width="2"/>
							<path stroke-linecap="round" stroke-width="2" d="M12 2v3m0 14v3M2 12h3m14 0h3" />
						</svg>
					{:else if item.icon === 'calendar'}
						<svg class={`w-4 h-4 ${active ? 'text-white' : 'text-slate-400 group-hover:text-slate-200'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
						</svg>
					{:else if item.icon === 'file-text'}
						<svg class={`w-4 h-4 ${active ? 'text-white' : 'text-slate-400 group-hover:text-slate-200'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
						</svg>
					{:else if item.icon === 'history'}
						<svg class={`w-4 h-4 ${active ? 'text-white' : 'text-slate-400 group-hover:text-slate-200'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					{:else if item.icon === 'smartphone'}
						<svg class={`w-4 h-4 ${active ? 'text-white' : 'text-slate-400 group-hover:text-slate-200'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<rect x="7" y="2" width="10" height="20" rx="2" stroke-width="2"/>
							<path stroke-linecap="round" stroke-width="2" d="M11 18h2" />
						</svg>
					{/if}

					<span>{item.label}</span>
				</div>

				<!-- Right Indicators / Badges -->
				{#if item.indicator === 'active-dot' && active}
					<span class="w-2 h-2 rounded-full bg-emerald-400 shadow-sm shadow-emerald-400/50"></span>
				{:else if item.badge}
					<span
						class={`text-[10px] px-1.5 py-0.5 rounded-full font-bold uppercase ${
							item.badge.color === 'green'
								? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
								: 'bg-amber-500/20 text-amber-400 border border-amber-500/30'
						}`}
					>
						{item.badge.text}
					</span>
				{/if}
			</a>
		{/each}

		<!-- Section Administrasi -->
		<div class="pt-4 pb-1">
			<div class="text-[10px] font-bold tracking-wider text-slate-400 uppercase px-3 mb-1">
				Administrasi
			</div>
			{#each adminItems as item}
				<a
					href={item.href}
					class="flex items-center justify-between px-3 py-2 rounded-lg text-slate-400 hover:text-slate-200 hover:bg-slate-800/50 transition-colors group"
				>
					<div class="flex items-center gap-3">
						{#if item.icon === 'building'}
							<svg class="w-4 h-4 text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
							</svg>
						{:else if item.icon === 'settings'}
							<svg class="w-4 h-4 text-slate-400 group-hover:text-slate-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
							</svg>
						{:else if item.icon === 'shield-check'}
							<svg class="w-4 h-4 text-slate-400 group-hover:text-slate-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
							</svg>
						{/if}
						<span class={item.ownerOnly ? 'text-slate-200 font-semibold' : ''}>{item.label}</span>
					</div>

					{#if item.badge}
						<span class="text-[9px] px-1.5 py-0.5 rounded-full font-bold uppercase bg-purple-500/20 text-purple-400 border border-purple-500/30">
							{item.badge.text}
						</span>
					{/if}
				</a>
			{/each}
		</div>
	</nav>

	<!-- Bottom User Profile Footer -->
	<div class="p-3 border-t border-slate-800/80 bg-[#090e1a]">
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-2.5 min-w-0">
				<!-- Avatar with initial or photo -->
				<div class="relative shrink-0">
					<div class="w-9 h-9 rounded-full bg-gradient-to-tr from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-xs shadow-md overflow-hidden">
						{#if auth.user?.name}
							<span>{auth.user.name.slice(0, 2).toUpperCase()}</span>
						{:else}
							<span>SR</span>
						{/if}
					</div>
					<span class="absolute bottom-0 right-0 w-2.5 h-2.5 rounded-full bg-emerald-500 ring-2 ring-[#090e1a]"></span>
				</div>
				<div class="min-w-0">
					<div class="text-xs font-semibold text-white truncate">
						{auth.user?.name || 'Siti Rahmadhani'}
					</div>
					<div class="text-[10px] text-slate-400 truncate">
						{auth.user?.role === 'OWNER' ? 'Owner / Founder' : auth.user?.role === 'ADMIN_COMPANY' ? 'HR Operations Head' : 'Employee'}
					</div>
				</div>
			</div>

			<!-- Logout Button -->
			<button
				type="button"
				class="p-1.5 text-slate-400 hover:text-rose-400 hover:bg-slate-800/60 rounded-lg transition-colors"
				title="Keluar Akun"
				onclick={() => auth.logout()}
				aria-label="Logout"
			>
				<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
				</svg>
			</button>
		</div>
	</div>
</aside>