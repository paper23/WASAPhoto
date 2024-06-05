<script>
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
                this.$router.push({ path: '/users/' + this.token });
            }
            catch(e) {
                this.errormsg = e.toString();
            }
        },
	},
	mounted() {
	}
}
</script>

<template>
	<h1>Carica un'immagine</h1>
    <form @submit.prevent="uploadPhoto" method="post" enctype="multipart/form-data">
        <label for="file">Scegli un'immagine:</label>
        <input type="file" id="file" name="file" accept="image/*" @change="handleFileUpload" required>
        <br><br>
        <input type="submit" value="Carica Immagine">
    </form>
	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
</style>