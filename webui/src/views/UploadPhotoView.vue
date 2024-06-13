<script>
import Navbar from '@/components/Navbar.vue';
import Toolbar from '@/components/Toolbar.vue';
export default {
	data: function() {
		return {
			errormsg: null,
			username: localStorage.getItem("username"),
			token: localStorage.getItem("token"),
			file: null,
		}
	},
	methods: {

		handleFileUpload(event) {
			this.file = event.target.files[0];
		},

        async uploadPhoto() {
			if (!this.file) {
                this.errormsg = "Please select a file";
                return;
			}
            try {
				let formData = new FormData();
                formData.append('file', this.file);

                let response = await this.$axios.post("/users/" + this.token + "/images/", formData, {
						headers: {
							Authorization: "Bearer " + this.token,
							'Content-Type': 'multipart/form-data'
						}})

				this.$router.push({ path: '/users/' + this.token, query: { username: this.username }});
            }
            catch(e) {
                this.errormsg = e.toString();
            }
        },
	},
	mounted() {
	},
	components: {
		Navbar,
		Toolbar
	},
}
</script>

<template>
	<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
		<h1 class="mb-4">Upload Photo</h1>
		<Toolbar />
	</div>
	<div>
		<Navbar />
		<div class="container d-flex justify-content-center align-items-center">
			<div class="text-center w-50">
				<form @submit.prevent="uploadPhoto" method="post" enctype="multipart/form-data" class="border p-4 rounded">
				<div class="form-group">
					<label for="file" class="form-label">Choose a photo:</label>
					<input type="file" id="file" name="file" accept="image/*" @change="handleFileUpload" class="form-control" required>
				</div>
				<br>
				<button type="submit" class="btn btn-primary">UPLOAD</button>
				</form>
			</div>
		</div>
	</div>
  </template>

<style>
</style>