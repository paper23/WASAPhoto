<script>
import Navbar from '@/components/Navbar.vue';
import Toolbar from '@/components/Toolbar.vue';
export default {
    data: function() {
        return {
            errormsg: null,
			username: localStorage.getItem("username"),
			token: localStorage.getItem("token"),

            bannedUsers: [
                {
                    idUser: 0,
                    username: "",
                }
            ],
        }
    },
    methods: {
        async getBannedList() {
            try {
            	let response = await this.$axios.get("/users/" + this.token + "/bans/", {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("token")
						}})
                this.bannedUsers = response.data
			}
			catch(e) {
				this.errormsg = e.toString();
			}
        }
    },

    components: {
        Navbar,
        Toolbar,
    },

    mounted() {
        this.getBannedList()
    }
}
</script>

<template>
    <div>
        <Navbar />
        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
            <h3>List of users you have banned</h3>
            <Toolbar />
        </div>
        <div class="row" v-if="this.bannedUsers[0].idUser != 0">
            <ul>
                <li v-for="banned in this.bannedUsers" :key="banned.idUser">
                    {{ banned.username }}
                    <button class="btn btn-success">Unban</button>
                </li>
            </ul>
        </div>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
</style>