import apiAuth from '@/api/auth.js'
import apiRole from '@/api/role.js'
import apiUser from '@/api/user.js'

import { basePermissions } from '@/settings'

export async function getUserInfo() {
  const res = await apiUser.getUser()
  const { id, username, profile, roles, currentRole } = res.data || {}
  return {
    id,
    username,
    avatar: profile?.avatar,
    nickName: profile?.nickName,
    gender: profile?.gender,
    address: profile?.address,
    email: profile?.email,
    roles,
    currentRole,
  }
}

export async function getPermissions() {
  let asyncPermissions = []
  try {
    const res = await apiRole.getRolePermissions()
    asyncPermissions = res?.data || []
  }
  catch (error) {
    console.error(error)
  }
  return basePermissions.concat(asyncPermissions)
}
