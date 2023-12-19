<script>
	/** @type {import('./$types').PageData*/
	export let data;
</script>

<main>
	<h1>Registers</h1>
	<table>
		<thead>
			<tr>
				<th>Tag</th>
				<th>Description</th>
				<th>Address</th>
				<th>Value</th>
				<th>Last Update</th>
			</tr>
		</thead>
		{#if data.data}
			{#each data.data as register}
				<tr>
					<td>{register.tag}</td>
					<td>{register.description}</td>
					<td>{register.address}</td>
					<td>{register.value}</td>
					<td>{register.last_update}</td>
				</tr>
			{/each}
		{/if}
	</table>
</main>
