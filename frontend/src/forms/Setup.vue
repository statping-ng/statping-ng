<template>
    <div class="container col-md-7 col-sm-12 mt-2 sm-container">
        <div class="col-12 col-md-6 offset-md-3 mb-4">
            <img alt="Statping-ng Setup" class="img-fluid mt-5 mt-md-0 logo" src="data:image/png;base64,PHN2ZyB3aWR0aD0iNTciIGhlaWdodD0iNTciIHZpZXdCb3g9IjAgMCA1NyA1NyIgZmlsbD0ibm9u
      ZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTIzLjg5Njgg
      MzAuM0MxOS4zNzIxIDMwLjMgMTUuNjk2OCAyNi42MDIzIDE1LjY5NjggMjIuMDVDMTUuNjk2OCAx
      Ny40OTc3IDE5LjM3MjEgMTMuOCAyMy44OTY4IDEzLjhDMjYuNDc3NyAxMy44IDI4Ljc2NDYgMTQu
      OTk5NyAzMC4yNjc0IDE2Ljg1NjhDMzIuMTYyMiAxNS4zMTIgMzQuMzUxMSAxMy41MzcxIDM2LjYz
      NzkgMTEuNjYzNUMzMy42MzIzIDcuOTMyOTcgMjkuMDQyMyA1LjU1IDIzLjg5NjggNS41NUMxNC44
      MzExIDUuNTUgNy40OTY4MyAxMi45NDU0IDcuNDk2ODMgMjIuMDVDNy40OTY4MyAzMS4xNzEgMTQu
      ODQ3NCAzOC41NSAyMy44OTY4IDM4LjU1QzI5LjA0MjMgMzguNTUgMzMuNjMyMyAzNi4xNjcgMzYu
      NjM3OSAzMi40MzY1TDMwLjI2NzQgMjcuMjQzMkMyOC43NjQ2IDI5LjEwMDMgMjYuNDc3NyAzMC4z
      IDIzLjg5NjggMzAuM1oiIGZpbGw9IiMwMDgyREUiLz4KPHBhdGggZD0iTTMwLjI2NzQgMjcuMjI2
      OEMzMS40MTA4IDI1LjgxMzQgMzIuMDk2OSAyNC4wMDU2IDMyLjA5NjkgMjIuMDMzNUMzMi4wOTY5
      IDIwLjA2MTQgMzEuNDEwOCAxOC4yNTM3IDMwLjI2NzQgMTYuODQwM0MyNi41MTA0IDE5Ljg5NzEg
      MjMuODk2OSAyMi4wMzM1IDIzLjg5NjkgMjIuMDMzNUwzMC4yNjc0IDI3LjIyNjhaIiBmaWxsPSIj
      MDBBQ0U4Ii8+CjxwYXRoIGQ9Ik00OC40OTY4IDIyLjA1QzQ4LjQ5NjggMTYuMTUwMSA0Ni40Mzg2
      IDEwLjcyNjggNDMuMDA4NCA2LjQ3MDMyQzQwLjkwMTIgOC4xOTU5MiAzOC43Mjg3IDkuOTcwODIg
      MzYuNjM3OSAxMS42NjM1QzM4LjkyNDggMTQuNTA2NyA0MC4yOTY4IDE4LjEwNTggNDAuMjk2OCAy
      Mi4wNUM0MC4yOTY4IDI1Ljk3NzggMzguOTI0OCAyOS41OTMzIDM2LjYzNzkgMzIuNDM2NUw0My4w
      MDg0IDM3LjYyOTdDNDYuNDU1IDMzLjM3MzIgNDguNDk2OCAyNy45NDk5IDQ4LjQ5NjggMjIuMDVa
      IiBmaWxsPSIjMDBBQ0U4Ii8+CjxwYXRoIGQ9Ik0zNi42NTQyIDExLjY2MzVDMzQuMzY3NCAxMy41
      MzcxIDMyLjE3ODUgMTUuMzEyIDMwLjI4MzcgMTYuODU2OEMzMS40MjcxIDE4LjI3MDEgMzIuMTEz
      MiAyMC4wNzc5IDMyLjExMzIgMjIuMDVDMzIuMTEzMiAyNC4wMjIxIDMxLjQyNzEgMjUuODI5OSAz
      MC4yODM3IDI3LjI0MzJMMzYuNjU0MiAzMi40MzY1QzM4Ljk0MTEgMjkuNTkzMyA0MC4zMTMyIDI1
      Ljk5NDIgNDAuMzEzMiAyMi4wNUM0MC4yOTY4IDE4LjEwNTggMzguOTQxMSAxNC40OTAyIDM2LjY1
      NDIgMTEuNjYzNVoiIGZpbGw9IiMwMDU0OTkiLz4KPHBhdGggZD0iTTIzLjkxOTQgNDguNDEyOUMy
      My41NTU5IDQ4LjczMzYgMjMuMTAzOSA0OC44ODUgMjIuNTg5OCA0OC45MDI4QzIyLjA2NjggNDgu
      ODg1IDIxLjU5NyA0OC43MTU3IDIxLjIzMzYgNDguMzg2MkMyMC44NjEzIDQ4LjAzODkgMjAuNjc1
      MSA0Ny41OTM1IDIwLjY4NCA0Ny4wNTAyQzIwLjY4NCA0Ni41MTU4IDIwLjg3OSA0Ni4wNjE2IDIx
      LjI0MjQgNDUuNzA1NEMyMS41OTcgNDUuMzU4IDIyLjA0OSA0NS4xNzk5IDIyLjU4MDkgNDUuMTYy
      MUMyMy4wODYxIDQ1LjE2MjEgMjMuNTIwNSA0NS4zMjI0IDIzLjkwMTYgNDUuNTk4NUwyNS4wNTQg
      NDQuNDIyOEMyNC42MDE5IDQ0LjA0ODggMjQuMDg3OCA0My43OTk0IDIzLjQ5MzkgNDMuNjgzNkMy
      My4yNjM0IDQzLjYzMDIgMjMuMDE1MiA0My42MTIzIDIyLjc2NyA0My42MDM0SDIyLjU4MDlDMjIu
      NTM2NiA0My42MDM0IDIyLjQ5MjMgNDMuNjAzNCAyMi40NDc5IDQzLjYxMjNIMjIuNDMwMkMyMi4x
      ODIgNDMuNjIxMyAyMS45MzM4IDQzLjY1NjkgMjEuNzAzMyA0My43MTAzQzIxLjA1NjMgNDMuODUy
      OCAyMC40OTc4IDQ0LjE0NjcgMjAuMDE5MiA0NC41OTIxQzE5LjMzNjYgNDUuMjQyMiAxOC45ODIx
      IDQ2LjA2MTYgMTguOTczMiA0Ny4wNjhDMTguOTY0MyA0OC4wNzQ1IDE5LjMxIDQ4Ljg4NSAxOS45
      OTI2IDQ5LjUxNzNDMjAuMjIzIDQ5LjczMTEgMjAuNDcxMiA0OS45MDAzIDIwLjc0NiA1MC4wNDI4
      QzIxLjI3NzkgNTAuMzEgMjEuODcxOCA1MC40NDM2IDIyLjUzNjYgNTAuNDM0N0gyMi41ODk4QzIz
      LjI1NDYgNTAuNDI1OCAyMy44MjE5IDUwLjI0NzYgMjQuMzUzNyA0OS45NzE1QzI0LjU5MyA0OS44
      NDY5IDI0LjgyMzUgNDkuNjg2NSAyNS4wMjc0IDQ5LjQ5OTVMMjMuOTE5NCA0OC40MTI5Wk01MC42
      MDkxIDQ4LjQxMjlDNTAuMjQ1NiA0OC43MzM2IDQ5Ljc5MzYgNDguODg1IDQ5LjI3OTUgNDguOTAy
      OEM0OC43NTY1IDQ4Ljg4NSA0OC4yODY3IDQ4LjcxNTcgNDcuOTIzMyA0OC4zODYyQzQ3LjU1MSA0
      OC4wMzg5IDQ3LjM2NDggNDcuNTkzNSA0Ny4zNzM3IDQ3LjA1MDJDNDcuMzczNyA0Ni41MTU4IDQ3
      LjU2ODcgNDYuMDYxNiA0Ny45MzIxIDQ1LjcwNTRDNDguMjg2NyA0NS4zNTggNDguNzM4OCA0NS4x
      Nzk5IDQ5LjI3MDYgNDUuMTYyMUM0OS43NzU4IDQ1LjE2MjEgNTAuMjEwMiA0NS4zMjI0IDUwLjU5
      MTMgNDUuNTk4NUw1MS43NDM3IDQ0LjQyMjhDNTEuMjkxNiA0NC4wNDg4IDUwLjc3NzUgNDMuNzk5
      NCA1MC4xODM2IDQzLjY4MzZDNDkuOTUzMSA0My42MzAyIDQ5LjcwNDkgNDMuNjEyMyA0OS40NTY3
      IDQzLjYwMzRINDkuMjcwNkM0OS4yMjYzIDQzLjYwMzQgNDkuMTgyIDQzLjYwMzQgNDkuMTM3NiA0
      My42MTIzSDQ5LjExOTlDNDguODcxNyA0My42MjEzIDQ4LjYyMzUgNDMuNjU2OSA0OC4zOTMxIDQz
      LjcxMDNDNDcuNzQ2IDQzLjg1MjggNDcuMTg3NSA0NC4xNDY3IDQ2LjcwODkgNDQuNTkyMUM0Ni4w
      MjYzIDQ1LjI0MjIgNDUuNjcxOCA0Ni4wNjE2IDQ1LjY2MjkgNDcuMDY4QzQ1LjY1NDEgNDguMDc0
      NSA0NS45OTk4IDQ4Ljg4NSA0Ni42ODIzIDQ5LjUxNzNDNDYuOTEyOCA0OS43MzExIDQ3LjE2MDkg
      NDkuOTAwMyA0Ny40MzU3IDUwLjA0MjhDNDcuOTY3NiA1MC4zMSA0OC41NjE1IDUwLjQ0MzYgNDku
      MjI2MyA1MC40MzQ3SDQ5LjI3OTVDNDkuOTQ0MyA1MC40MjU4IDUwLjUxMTYgNTAuMjQ3NiA1MS4w
      NDM0IDQ5Ljk3MTVDNTEuMjgyNyA0OS44NDY5IDUxLjUxMzIgNDkuNjg2NSA1MS43MTcxIDQ5LjQ5
      OTVMNTAuNjA5MSA0OC40MTI5Wk0yOC44NjU1IDUwLjQ2MTRDMjcuODYzOSA1MC40NjE0IDI3LjAz
      MDcgNTAuMTMxOSAyNi4zNjU4IDQ5LjQ1NUMyNS43MDEgNDguNzc4MSAyNS4zNzMxIDQ3Ljk2NzYg
      MjUuMzczMSA0Ny4wMDU3QzI1LjM3MzEgNDYuMDQzOCAyNS43MDEgNDUuMjMzMyAyNi4zNjU4IDQ0
      LjU1NjRDMjcuMDMwNyA0My44Nzk1IDI3Ljg2MzkgNDMuNTUgMjguODY1NSA0My41NUMyOS44Njcy
      IDQzLjU1IDMwLjcwMDQgNDMuODc5NSAzMS4zNjUyIDQ0LjU1NjRDMzIuMDMgNDUuMjMzMyAzMi4z
      NTggNDYuMDQzOCAzMi4zNTggNDcuMDA1N0MzMi4zNTggNDcuOTY3NiAzMi4wMyA0OC43NzgxIDMx
      LjM2NTIgNDkuNDU1QzMwLjcwMDQgNTAuMTMxOSAyOS44NjcyIDUwLjQ2MTQgMjguODY1NSA1MC40
      NjE0Wk0yOC44NjU1IDQ4Ljg4NUMyOS4zNzA4IDQ4Ljg4NSAyOS43Njk2IDQ4LjcwNjggMzAuMTA2
      NSA0OC4zNTA2QzMwLjQzNDUgNDguMDAzMiAzMC42MTE3IDQ3LjU0MDEgMzAuNjExNyA0Ny4wMDU3
      QzMwLjYxMTcgNDYuNDYyNCAzMC40NDMzIDQ2LjAwODIgMzAuMTA2NSA0NS42NjA4QzI5Ljc3ODUg
      NDUuMzEzNSAyOS4zNzA4IDQ1LjEyNjQgMjguODY1NSA0NS4xMjY0QzI4LjM2MDMgNDUuMTI2NCAy
      Ny45NTI1IDQ1LjMwNDYgMjcuNjE1NyA0NS42NjA4QzI3LjI4NzcgNDYuMDA4MiAyNy4xMjgyIDQ2
      LjQ3MTMgMjcuMTI4MiA0Ny4wMDU3QzI3LjEyODIgNDcuNTQ5IDI3LjI4NzcgNDguMDAzMiAyNy42
      MTU3IDQ4LjM1MDZDMjcuOTQzNyA0OC43MDY4IDI4LjM2MDMgNDguODg1IDI4Ljg2NTUgNDguODg1
      Wk0xNC40NjE0IDUwLjQ2MTRDMTMuNDU5OCA1MC40NjE0IDEyLjYyNjUgNTAuMTMxOSAxMS45NjE3
      IDQ5LjQ1NUMxMS4yOTY5IDQ4Ljc3ODEgMTAuOTY5IDQ3Ljk2NzYgMTAuOTY5IDQ3LjAwNTdDMTAu
      OTY5IDQ2LjA0MzggMTEuMjk2OSA0NS4yMzMzIDExLjk2MTcgNDQuNTU2NEMxMi42MjY1IDQzLjg3
      OTUgMTMuNDU5OCA0My41NSAxNC40NjE0IDQzLjU1QzE1LjQ2MyA0My41NSAxNi4yOTYzIDQzLjg3
      OTUgMTYuOTYxMSA0NC41NTY0QzE3LjYyNTkgNDUuMjMzMyAxNy45NTM4IDQ2LjA0MzggMTcuOTUz
      OCA0Ny4wMDU3QzE3Ljk1MzggNDcuOTY3NiAxNy42MjU5IDQ4Ljc3ODEgMTYuOTYxMSA0OS40NTVD
      MTYuMjk2MyA1MC4xMzE5IDE1LjQ2MyA1MC40NjE0IDE0LjQ2MTQgNTAuNDYxNFpNMTQuNDYxNCA0
      OC44ODVDMTQuOTY2NyA0OC44ODUgMTUuMzY1NSA0OC43MDY4IDE1LjcwMjQgNDguMzUwNkMxNi4w
      MzAzIDQ4LjAwMzIgMTYuMjA3NiA0Ny41NDAxIDE2LjIwNzYgNDcuMDA1N0MxNi4yMDc2IDQ2LjQ2
      MjQgMTYuMDM5MiA0Ni4wMDgyIDE1LjcwMjQgNDUuNjYwOEMxNS4zNzQ0IDQ1LjMxMzUgMTQuOTY2
      NyA0NS4xMjY0IDE0LjQ2MTQgNDUuMTI2NEMxMy45NTYyIDQ1LjEyNjQgMTMuNTQ4NCA0NS4zMDQ2
      IDEzLjIxMTYgNDUuNjYwOEMxMi44ODM2IDQ2LjAwODIgMTIuNzI0IDQ2LjQ3MTMgMTIuNzI0IDQ3
      LjAwNTdDMTIuNzI0IDQ3LjU0OSAxMi44ODM2IDQ4LjAwMzIgMTMuMjExNiA0OC4zNTA2QzEzLjU0
      ODQgNDguNzA2OCAxMy45NjUgNDguODg1IDE0LjQ2MTQgNDguODg1Wk0zOC40NzQyIDQzLjY5MjVI
      MzIuNzkyM1Y0NS4xOTc3SDM0LjQzMjFDMzQuNjI3MiA0NS4xOTc3IDM0Ljc3NzggNDUuMzU4IDM0
      Ljc3NzggNDUuNTQ1VjUwLjM0NTZIMzYuNDc5N1Y0NS41NDVDMzYuNDc5NyA0NS4zNDkxIDM2LjYz
      OTMgNDUuMTk3NyAzNi44MjU0IDQ1LjE5NzdIMzguNDc0MlY0My42OTI1Wk00NC43MjMzIDQ4Ljg0
      OTNINDEuNzUzOUM0MS41NTg5IDQ4Ljg0OTMgNDEuNDA4MiA0OC42ODkgNDEuNDA4MiA0OC41MDJW
      NDcuNzAwNEg0My45NDMzVjQ2LjMwMjFINDEuNDA4MlY0NS41NjI5QzQxLjQwODIgNDUuMzY2OSA0
      MS41Njc3IDQ1LjIxNTUgNDEuNzUzOSA0NS4yMTU1SDQ0LjY0MzZWNDMuNzE5MkgzOS42OTc0VjUw
      LjM3MjNINDQuNzE0NVY0OC44NDkzSDQ0LjcyMzNaTTcuODkzMTMgNDYuMzU1NUM2LjU3MjM5IDQ2
      LjEwNjIgNi40MjE3IDQ1Ljg0NzkgNi40Mzk0MiA0NS41NzE4QzYuNDY2MDIgNDUuMjY4OSA2Ljgy
      OTQ0IDQ1LjExNzUgNy4zMzQ3IDQ1LjA5OTdDNy44NzU0IDQ1LjA4MTkgOC41MTM2MiA0NS4xNzk5
      IDkuMTc4NDIgNDUuNTM2MUwxMC4yMzMyIDQ0LjM4NzJDOS40MTc3NSA0My44NjE3IDguMzcxNzkg
      NDMuNTg1NiA3LjMwODEgNDMuNjIxM0M1LjcyMTQ0IDQzLjY3NDcgNC42OTMyIDQ0LjU1NjQgNC43
      Mzc1MiA0NS43NTg4VjQ1Ljc3NjZDNC43NzI5OCA0Ni45MTY2IDUuNjg1OTggNDcuNDMzMiA3LjE3
      NTE0IDQ3LjcxODJDOC40NzgxNiA0Ny45Njc2IDguNjM3NzEgNDguMTkwMyA4LjY0NjU4IDQ4LjQy
      MThWNDguNDM5NkM4LjY1NTQ0IDQ4Ljc0MjUgOC4xMzI0NiA0OC45NDczIDcuNTQ3NDMgNDguOTY1
      MUM2LjgxMTcyIDQ4Ljk5MTggNS45NjA3NiA0OC44MDQ4IDUuMjg3MSA0OC4yOTcxSDUuMjc4MjNM
      NC4yNSA0OS40MzcyTDQuMzM4NjQgNDkuNDk5NUM1LjIzMzkxIDUwLjE0OTcgNi4zNzczOCA1MC40
      NzkyIDcuNTc0MDMgNTAuNDQzNkM5LjMxMTM4IDUwLjM4MTIgMTAuMzgzOSA0OS40MjgzIDEwLjMz
      OTYgNDguMjQzN1Y0OC4yMjU5QzEwLjI5NTMgNDcuMTAzNyA5LjM2NDU3IDQ2LjY0MDUgNy44OTMx
      MyA0Ni4zNTU1WiIgZmlsbD0iYmxhY2siLz4KPC9zdmc+Cg==" />
        </div>

        <div class="col-12">
            <form @submit.prevent="saveSetup">
                <div class="row">
                    <div class="col-12 col-md-6">
                        <div class="form-group">
                            <label class="text-capitalize">{{ $t('language') }}</label>
                            <select @change="changeLanguages" v-model="setup.language" id="language" class="form-control">
                                <option value="en">English</option>
                                <option value="es">Spanish</option>
                                <option value="fr">French</option>
                                <option value="ru">Russian</option>
                                <option value="de">German</option>
                                <option value="cs">Czech</option>
                                <option value="ja">Japanese</option>
                                <option value="ko">Korean</option>
                                <option value="it">Italian</option>
                                <option value="zh">Chinese</option>
                                <option value="sv">Swedish</option>
                            </select>
                        </div>
                        <div class="form-group">
                            <label class="text-capitalize">{{ $t('db_connection') }}</label>
                            <select @change="canSubmit" v-model="setup.db_connection" id="db_connection" class="form-control">
                                <option value="sqlite">SQLite</option>
                                <option value="postgres">Postgres</option>
                                <option value="mysql">MySQL</option>
                            </select>
                        </div>
                        <div class="row">
                            <div class="col-7 col-md-6">
                                <div v-if="setup.db_connection !== 'sqlite'" class="form-group">
                                    <label class="text-capitalize">{{ $t('db_host') }}</label>
                                    <input @keyup="canSubmit" v-model="setup.db_host" id="db_host" type="text" class="form-control" placeholder="localhost">
                                </div>
                            </div>
                            <div class="col-5 col-md-6">
                                <div v-if="setup.db_connection !== 'sqlite'" class="form-group">
                                    <label class="text-capitalize">{{ $t('db_port') }}</label>
                                    <input @keyup="canSubmit" v-model.number="setup.db_port" id="db_port" type="number" class="form-control" placeholder="5432">
                                </div>
                            </div>
                        </div>
                        <div v-if="setup.db_connection !== 'sqlite'" class="form-group">
                            <label class="text-capitalize">{{ $t('db_username') }}</label>
                            <input @keyup="canSubmit" v-model="setup.db_user" id="db_user" type="text" class="form-control" placeholder="root">
                        </div>
                        <div v-if="setup.db_connection !== 'sqlite'" class="form-group">
                            <label for="db_password" class="text-capitalize">{{ $t('db_password') }}</label>
                            <input @keyup="canSubmit" v-model="setup.db_password" id="db_password" type="password" class="form-control" placeholder="password123">
                        </div>
                        <div v-if="setup.db_connection !== 'sqlite'" class="form-group">
                            <label for="db_database" class="text-capitalize">{{ $t('db_database') }}</label>
                            <input @keyup="canSubmit" v-model="setup.db_database" id="db_database" type="text" class="form-control" placeholder="Database name">
                        </div>

                        <div class="form-group mt-3">
                            <div class="row">
                                <div class="col-9">
                                    <span class="text-left text-capitalize">{{ $t('send_reports') }}</span>
                                </div>
                                <div class="col-3 text-right">
                                    <span @click="setup.send_reports = !!setup.send_reports" class="switch">
                                      <input v-model="setup.send_reports" type="checkbox" name="send_reports" class="switch" id="send_reports" :checked="setup.send_reports">
                                      <label for="send_reports"></label>
                                    </span>
                                </div>
                            </div>
                        </div>

                    </div>

                    <div class="col-12 col-md-6">

                        <div class="form-group">
                            <label class="text-capitalize">{{ $t('project_name') }}</label>
                            <input @keyup="canSubmit" v-model="setup.project" id="project" type="text" class="form-control" placeholder="Work Servers" required>
                        </div>

                        <div class="form-group">
                            <label class="text-capitalize">{{ $t('description') }}</label>
                            <input @keyup="canSubmit" v-model="setup.description" id="description" type="text" class="form-control" placeholder="Monitors all of my work services">
                        </div>

                        <div class="form-group">
                            <label class="text-capitalize" for="domain">{{ $t('domain') }}</label>
                            <input @keyup="canSubmit" v-model="setup.domain" type="text" class="form-control" id="domain" required>
                        </div>

                        <div class="form-group">
                            <label class="text-capitalize">{{ $t('username') }}</label>
                            <input @keyup="canSubmit" v-model="setup.username" id="username" type="text" class="form-control" placeholder="admin" required>
                        </div>

                        <div class="form-group">
                            <label class="text-capitalize">{{ $t('password') }}</label>
                            <input @keyup="canSubmit" v-model="setup.password" id="password" type="password" class="form-control" placeholder="password" required>
                        </div>

                        <div class="form-group">
                            <label class="text-capitalize">{{ $t('confirm_password') }}</label>
                            <input @keyup="canSubmit" v-model="setup.confirm_password" id="password_confirm" type="password" class="form-control" placeholder="password" required>
                            <span v-if="passnomatch" class="small text-danger">Both passwords should match</span>
                        </div>

                        <div class="form-group">
                            <div class="row">
                                <div class="col-8">
                                    <label class="text-capitalize">{{ $t('email') }}</label>
                                    <input @keyup="canSubmit" v-model="setup.email" id="email" type="text" class="form-control" placeholder="myemail@domain.com">
                                </div>
                                <div class="col-4 text-right">
                                    <label class="d-none d-sm-block text-capitalize text-capitalize">{{ $t('newsletter') }}</label>
                                    <span @click="setup.newsletter = !!setup.newsletter" class="switch">
                                      <input v-model="setup.newsletter" type="checkbox" name="send_newsletter" class="switch" id="send_newsletter" :checked="setup.newsletter">
                                      <label for="send_newsletter"></label>
                                    </span>
                                </div>
                            </div>
                            <small>{{ $t('newsletter_note') }}</small>
                        </div>
                    </div>

                    <div v-if="error" class="col-12 alert alert-danger">
                        {{error}}
                    </div>

                    <div class="col-12">
                        <button @click.prevent="saveSetup" v-bind:disabled="disabled || loading" type="submit" class="btn btn-primary btn-block" :class="{'btn-primary': !loading, 'btn-default': loading}">
                            <font-awesome-icon v-if="loading" icon="circle-notch" class="mr-2" spin/>{{loading ? $t('loading') : $t('save_settings')}}
                        </button>
                    </div>
                </div>
            </form>

        </div>
    </div>
</template>

<script>
  import Api from "../API";

  export default {
  name: 'Setup',
  data () {
    return {
      error: null,
      loading: false,
      disabled: true,
      passnomatch: false,
      setup: {
        language: "en",
        db_connection: "sqlite",
        db_host: "",
        db_port: "",
        db_user: "",
        db_password: "",
        db_database: "",
        project: "",
        description: "",
        domain: "",
        username: "",
        password: "",
        confirm_password: "",
        sample_data: true,
        newsletter: true,
        send_reports: true,
        email: "",
      }
    }
  },
  async created() {
    const core = await Api.core()
    if (core.setup) {
        if (!this.$store.getters.hasPublicData) {
            await this.$store.dispatch('loadRequired')
        }
        this.$router.push('/')
    }
  },
  mounted() {
    this.changeLanguages()
    this.setup.domain = window.location.protocol + "//" + window.location.hostname + (window.location.port ? ":"+window.location.port : "")
  },
  methods: {
    changeLanguages() {
      this.$i18n.locale = this.setup.language
    },
      canSubmit() {
          this.error = null
          const s = this.setup
        if (s.confirm_password.length > 0 && s.confirm_password !== s.password) {
          this.passnomatch = true
        } else {
          this.passnomatch = false
        }
          if (s.db_connection !== 'sqlite') {
              if (!s.db_host || !s.db_port || !s.db_user || !s.db_password || !s.db_database) {
                  this.disabled = true
                  return
              }
          }
          if (!s.project || !s.domain || !s.username || !s.password || !s.confirm_password) {
              this.disabled = true
              return
          }
          if (s.password !== s.confirm_password) {
              this.disabled = true
              return
          }
          this.disabled = false
      },
    async saveSetup() {
      this.loading = true
      let resp
      try {
        resp = await Api.setup_save(this.setup)
      } catch(e) {
        resp = {status: 'error', error: e.response.data.error}
      }
      if (resp.status === 'error') {
        this.error = resp.error
        this.loading = false
        return
      }

      await this.$store.dispatch('loadCore')
      await this.$store.dispatch('loadRequired')

      this.loading = false
      this.$router.push('/')
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
