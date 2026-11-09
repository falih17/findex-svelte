<script>
	import '../app.css';
	import Navbar from '$lib/components/Navbar.svelte';
	let { children } = $props();
</script>

<div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col antialiased">
	<Navbar />
	<main class="flex-1 w-full max-w-5xl mx-auto p-4 sm:p-6">
		{@render children()}
	</main>
	<footer class="py-4 text-center text-xs text-slate-600 border-t border-slate-900">
		&copy; 2026 Findex Multi-Tenant SaaS Attendance System. PostgreSQL PostGIS & pgvector inside.
	</footer>
</div>
