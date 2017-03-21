var homeView = Vue.extend({
    template: `
        <div>
            <div class="page-wrapper">
                <div class="page-wrapper-row">
                    <div class="page-wrapper-top">
                        <div class="page-header">
                            <div class="page-header-top">
                                <div class="container">
                                    <a href="javascript:;" class="menu-toggler"></a>
                                    <div class="top-menu">
                                        <ul class="nav navbar-nav pull-right">
                                            <li class="droddown dropdown-separator">
                                                <span class="separator"></span>
                                            </li>
                                            <li class="dropdown dropdown-user dropdown-dark">
                                                <a href="javascript:;" class="dropdown-toggle" data-toggle="dropdown" data-hover="dropdown" data-close-others="true">
                                                    <span class="username username-hide-mobile">Relief Center Lenia Word</span>
                                                </a>
                                            </li>
                                            <a href="/logout">Log out</a>
                                        </ul>
                                    </div>
                                </div>
                            </div>
                            <div class="page-header-menu">
                                <div class="container">
                                    <form id="searchForm" class="search-form" action="/search" method="GET">
                                        <div class="input-group">
                                            <input type="text" class="form-control" placeholder="Search" name="query">
                                            <span class="input-group-btn">
                                                <a href="javascript:;" class="btn submit">
                                                    <i class="icon-magnifier"></i>
                                                </a>
                                            </span>
                                        </div>
                                    </form>
                                    <div class="hor-menu">
                                        <ul class="nav navbar-nav">
                                            <li aria-haspopup="true" class="menu-dropdown classic-menu-dropdown active">
                                                <a href="javascript:;"> Dashboard
                                                    <span class="arrow"></span>
                                                </a>
                                            </li>
                                            <li aria-haspopup="true" class="menu-dropdown mega-menu-dropdown  ">
                                                <a href="javascript:;"> Appointments
                                                    <span class="arrow"></span>
                                                </a>
                                            </li>
                                            <li aria-haspopup="true" class="menu-dropdown classic-menu-dropdown ">
                                                <a href="javascript:;"> Templates
                                                    <span class="arrow"></span>
                                                </a>
                                            </li>
                                        </ul>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="page-wrapper-row full-height">
                    <div class="page-wrapper-middle">
                        <div class="page-container">
                            <div class="page-content-wrapper">
                                <div class="page-head">
                                    <div class="container">
                                        <div class="page-title">
                                            <h1>{{ customer.Firstname }} {{ customer.Lastname }}
                                                <small>{{ customer.Address }} {{ customer.Phone }}
                                                </small>
                                            </h1>
                                        </div>
                                    </div>
                                </div>
                                <div class="page-content">
                                    <div class="container">
                                        <div class="page-content-inner">
                                            <div class="mt-content-body">
                                                <div class="row">
                                                    <div class="col*-12">
                                                        <div class="portlet light ">
                                                            <div class="portlet-title">
                                                                <div class="caption caption-md">
                                                                    <i class="icon-bar-chart font-dark hide"></i>
                                                                    <span class="caption-subject font-green-steel uppercase bold">Communication Logs</span>
                                                                </div>
                                                            </div>
                                                            <div class="portlet-body">
                                                                <table class="table table-hover" id="grid-template">
                                                                    <thead>
                                                                        <tr>
                                                                            <th> Date </th>
                                                                            <th> Direction </th>
                                                                            <th> Type </th>
                                                                        </tr>
                                                                    </thead>
                                                                    <tbody>
                                                                        <tr v-for="log in commLogs">
                                                                            <td>{{ log.Date }}</td>
                                                                            <td>{{ log.Direction }}</td>
                                                                            <td>{{ log.CommType }}</td>
                                                                        </tr>
                                                                    </tbody>
                                                                </table>
                                                            </div>
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="page-wrapper-row">
                    <div class="page-wrapper-bottom">
                    </div>
                </div>
            </div>
        </div>
    `,
    data: function() {
        return {
            customer: "",
            commLogs: []
        }
    },
    created: function() {
        console.log("Home object created.")

        this.$http.get('/dashboard').then(function(response) {
            this.customer = response.data ? response.data : ""
        })

        this.$http.get('/communication?pid=1234').then(function(response) {
            this.commLogs = response.data.CommLog ? response.data.CommLog : []
        })
    },
    methods: {
        logOut() {
            console.log("Logout.");
        },
        onSubmit() {
            console.log("Submit search form.")
        }
    }
});

var router = new VueRouter({
    mode: 'history',
    routes: [
        { path: '/', component: homeView },
        { path: '/home', component: homeView }
    ]   
});

router.beforeEach((to, from, next) => {
    console.log(`Navigation from ${from.path} to ${to.path}`)
    next()
})

var app = new Vue({
    router,
    template: `<router-view class="view"></router-view>`,
    created: function() {
         console.log("Vue object created.")
    }
}).$mount('#app');