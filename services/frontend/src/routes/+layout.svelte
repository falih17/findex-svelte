<script lang="ts">
	import '../app.css';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import TopHeader from '$lib/components/TopHeader.svelte';

	let { children } = $props();
	let sidebarOpen = $state(false);
</script>

<div class="min-h-screen bg-[#f8fafc] text-slate-800 flex antialiased">
	<!-- Left Navigation Sidebar -->
	<Sidebar bind:open={sidebarOpen} />

	<!-- Main Body Wrapper -->
	<div class="flex-1 lg:pl-64 flex flex-col min-h-screen min-w-0">
		<!-- Top App Bar -->
		<TopHeader onToggleSidebar={() => sidebarOpen = !sidebarOpen} />

		<!-- Page Content Canvas -->
		<main class="flex-1 p-4 sm:p-6 w-full max-w-7xl mx-auto space-y-6">
			{@render children()}
		</main>

		<!-- Subdued SaaS Footer -->
		<footer class="py-4 px-6 text-center text-[11px] text-slate-400 border-t border-slate-200/70 bg-white/50">
			&copy; 2026 AbsenPro / Findex Cloud SaaS Platform. PostGIS Geofencing & pgvector Facial Biometrics.
		</footer>
	</div>
</div>