<template>
  <div class="app-container">
    <!--dialog-->
    <el-dialog :visible.sync="dialogVisible" :title="$t('Edit User')">
      <el-form ref="form" :model="userForm" label-width="80px" :rules="rules" inline-message>
        <el-form-item prop="username" :label="$t('Username')" required>
          <el-input v-model="userForm.username" :placeholder="$t('Username')" :disabled="!isAdd"></el-input>
        </el-form-item>
        <el-form-item prop="password" :label="$t('Password')" required>
          <el-input type="password" v-model="userForm.password" :placeholder="$t('Password')"></el-input>
        </el-form-item>
        <el-form-item prop="role" :label="$t('Role')" required>
          <el-select v-model="userForm.role" :placeholder="$t('Role')">
            <el-option value="admin" :label="$t('admin')"></el-option>
            <el-option value="normal" :label="$t('normal')"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item prop="email" :label="$t('Email')">
          <el-input v-model="userForm.email" :placeholder="$t('Email')"/>
        </el-form-item>
      </el-form>
      <template slot="footer">
        <el-button size="small" @click="dialogVisible=false">{{$t('Cancel')}}</el-button>
        <el-button type="primary" size="small" @click="onConfirm">{{$t('Confirm')}}</el-button>
      </template>
    </el-dialog>
    <!--./dialog-->

    <el-card>
      <div class="filter">
        <div class="left"></div>
        <div class="right">
          <el-button type="success" icon="el-icon-plus" size="small" @click="onClickAddUser">添加用户</el-button>
        </div>
      </div>
      <!--table-->
      <el-table
        :data="userList"
        :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
        border
      >
        <el-table-column
          width="120px"
          :label="$t('Username')"
          prop="username"
        >
        </el-table-column>
        <el-table-column
          width="150px"
          :label="$t('Role')"
        >
          <template slot-scope="scope">
            <el-tag v-if="scope.row.role === 'admin'" type="primary">
              {{ $t(scope.row.role) }}
            </el-tag>
            <el-tag v-else type="warning">
              {{ $t(scope.row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          width="150px"
          :label="$t('Create Time')"
        >
          <template slot-scope="scope">
            {{getTime(scope.row.create_ts)}}
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('Action')"
          fixed="right"
        >
          <template slot-scope="scope">
            <el-button icon="el-icon-edit" type="warning" size="mini" @click="onEdit(scope.row)"></el-button>
            <el-button icon="el-icon-delete" type="danger" size="mini" @click="onRemove(scope.row)"></el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          @current-change="onPageChange"
          @size-change="onPageChange"
          :current-page.sync="pageNum"
          :page-sizes="[10, 20, 50, 100]"
          :page-size.sync="pageSize"
          layout="sizes, prev, pager, next"
          :total="totalCount">
        </el-pagination>
      </div>
      <!--./table-->
    </el-card>
  </div>
</template>

<script>
import {
  mapState
} from 'vuex'
import dayjs from 'dayjs'

export default {
  name: 'UserList',
  data () {
    const validatePass = (rule, value, callback) => {
      if (!value) return callback()
      if (value.length < 5) {
        callback(new Error(this.$t('Password length should be no shorter than 5')))
      } else {
        callback()
      }
    }
    const validateEmail = (rule, value, callback) => {
      if (!value) return callback()
      if (!value.match(/.+@.+/i)) {
        callback(new Error(this.$t('Email format invalid')))
      } else {
        callback()
      }
    }
    return {
      dialogVisible: false,
      isAdd: false,
      rules: {
        password: [{ validator: validatePass }],
        email: [{ validator: validateEmail }]
      }
    }
  },
  computed: {
    ...mapState('user', [
      'userList',
      'userForm',
      'totalCount'
    ]),
    pageSize: {
      get () {
        return this.$store.state.user.pageSize
      },
      set (value) {
        this.$store.commit('user/SET_PAGE_SIZE', value)
      }
    },
    pageNum: {
      get () {
        return this.$store.state.user.pageNum
      },
      set (value) {
        this.$store.commit('user/SET_PAGE_NUM', value)
      }
    }
  },
  methods: {
    onPageChange () {
      this.$store.dispatch('user/getUserList')
    },
    getTime (ts) {
      return dayjs(ts).format('YYYY-MM-DD HH:mm:ss')
    },
    onEdit (row) {
      this.isAdd = false
      this.$store.commit('user/SET_USER_FORM', row)
      this.dialogVisible = true
    },
    onRemove (row) {
      this.$confirm(this.$t('Are you sure to delete this user?'), this.$t('Notification'), {
        confirmButtonText: this.$t('Confirm'),
        cancelButtonText: this.$t('Cancel'),
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('user/deleteUser', row._id)
          .then(() => {
            this.$message({
              type: 'success',
              message: this.$t('Deleted successfully')
            })
          })
          .then(() => {
            this.$store.dispatch('user/getUserList')
          })
        this.$st.sendEv('用户列表', '删除用户')
      })
      // this.$store.commit('user/SET_USER_FORM', row)
    },
    onConfirm () {
      this.$refs.form.validate(valid => {
        if (!valid) return
        if (this.isAdd) {
          // 添加用户
          this.$store.dispatch('user/addUser')
            .then(() => {
              this.$message({
                type: 'success',
                message: this.$t('Saved successfully')
              })
              this.dialogVisible = false
              this.$st.sendEv('用户列表', '添加用户')
            })
            .then(() => {
              this.$store.dispatch('user/getUserList')
            })
        } else {
          // 编辑用户
          this.$store.dispatch('user/editUser')
            .then(() => {
              this.$message({
                type: 'success',
                message: this.$t('Saved successfully')
              })
              this.dialogVisible = false
              this.$st.sendEv('用户列表', '编辑用户')
            })
        }
      })
    },
    onClickAddUser () {
      this.isAdd = true
      this.$store.commit('user/SET_USER_FORM', {})
      this.dialogVisible = true
    },
    onValidateEmail (value) {
    }
  },
  created () {
    this.$store.dispatch('user/getUserList')
  }
}
</script>

<style scoped>
  .filter {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;

  .filter-search {
    width: 240px;
  }

  .right {

  .btn {
    margin-left: 10px;
  }

  }
  }

  .el-table {
    border-radius: 5px;
  }

  .el-table .el-button {
    padding: 7px;
  }
</style>
