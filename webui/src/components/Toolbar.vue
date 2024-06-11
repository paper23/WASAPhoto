<script>
    export default {
        data: function() {
            return {
                errormsg: null,
                username: localStorage.getItem("username"),
                token: localStorage.getItem("token"),
            }
        },
        methods: {
            async doLogout() {
                try {
                    localStorage.removeItem("token")
                    localStorage.removeItem("username")
                    this.$router.push({ path: '/' });
                }
                catch(e) {
                    this.errormsg = e.toString();
                }
            },
        }
    }
</script>
<template>
    <div class="btn-toolbar mb-2 mb-md-0">
		<div class="btn-group me-2">
			<button class="btn btn-danger" type="button" @click="doLogout">Logout</button>
		</div>
	</div>
</template>