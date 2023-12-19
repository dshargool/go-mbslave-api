/** @type {import('./$types').PageLoad} */
export async function load({ fetch, params }) {
	const resp = await fetch('http://localhost:8081/all_registers');
	if (resp.ok) {
		const data = await resp.json();
		console.log(data);
		return { data: data };
	} else {
		console.error('Failed to fetch data from PI');
		return {
			status: resp.status,
			error: new Error('Failed to fetch data')
		};
	}
}
