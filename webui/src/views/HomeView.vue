<script>
import Navbar from '@/components/Navbar.vue';
import Toolbar from '@/components/Toolbar.vue';
export default {
	data: function() {
		return {
			errormsg: null,
			username: localStorage.getItem("username"),
			token: localStorage.getItem("token"),
		}
	},
	
	components: {
		Navbar,
		Toolbar,
	},
}
</script>

<template>
	<div>
		<Navbar />
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Welcome back {{this.username }}</h1>
			<Toolbar />
		</div>
	</div>
</template>

<style>
</style>
