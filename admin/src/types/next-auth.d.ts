import NextAuth from 'next-auth'

declare module 'next-auth' {
  interface User {
    role: string
    accessToken: string
    refreshToken: string
    image?: string
  }

  interface Session {
    accessToken: string
    refreshToken: string
    user: {
      id: string
      email: string
      name: string
      role: string
      image?: string
    }
  }
}

declare module 'next-auth/jwt' {
  interface JWT {
    accessToken: string
    refreshToken: string
    role: string
  }
}
