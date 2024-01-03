<script>
	/** @type {import('./$types').PageData}*/
	import { onMount } from 'svelte';
    import { page } from '$app/stores';

	async function getApiStatus() {
		try {
			const response = await fetch($page.url.origin + '/api/healthcheck');

			if (response.ok) {
				// Request succeeded (status in the range of 200-299)
				return true;
			} else {
				// Request failed (status outside the range of 200-299)
				return false;
			}
		} catch (error) {
			// Network error or other issues
			console.error('Error occurred:', error);
			return false;
		}
	}

	async function updateStatus() {
		const api_status = await getApiStatus();
		const elem = document.getElementById('api-status');
		if (elem == null) {
			return;
		}

		if (api_status) {
			if (elem.classList.contains('badge-error')) {
				elem.classList.remove('badge-error');
			}
			elem.classList.add('badge-success');
		} else {
			if (elem.classList.contains('badge-success')) {
				elem.classList.remove('badge-success');
			}
			elem.classList.add('badge-error');
		}
		return;
	}

	onMount(async () => {
		setInterval(() => {
			updateStatus();
		}, 10000);
	});
</script>

<div class="flex flex-col min-h-screen justify-center gap-9 px-9 items-center">
	<h1>Modbus API</h1>
	<h3>Status</h3>
	<div class="indicator">
		<span
			use:updateStatus
			id="api-status"
			class="indicator-item indicator-middle indicator-center badge badge-lg"
		></span>
	</div>
</div>
