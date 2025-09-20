import NextAuth from 'next-auth'
import CredentialsProvider from 'next-auth/providers/credentials'
import { AuthOptions } from 'next-auth'

const authOptions: AuthOptions = {
  debug: true,
  providers: [
    CredentialsProvider({
      name: 'credentials',
      credentials: {
        email: { label: 'Email', type: 'email' },
        password: { label: 'Password', type: 'password' }
      },
      async authorize(credentials) {
        if (!credentials?.email || !credentials?.password) {
          console.warn('❌ Missing credentials')
          return null
        }

        const apiBaseUrl =
          process.env.NEXT_INTERNAL_API_URL || process.env.NEXT_PUBLIC_API_URL

        if (!apiBaseUrl) {
          console.error('❌ NEXT_INTERNAL_API_URL/NEXT_PUBLIC_API_URL not configured')
          throw new Error('API URL not configured')
        }

        const loginUrl = `${apiBaseUrl.replace(/\/$/, '')}/auth/login`

        try {
          const response = await fetch(loginUrl, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              email: credentials.email,
              password: credentials.password,
            }),
            cache: 'no-store',
          })

          const data = await response.json().catch(() => null)

          if (!response.ok) {
            const message = data?.message || 'Invalid credentials'
            console.warn('❌ Login failed:', message)
            return null
          }

          if (!data?.user || !data?.access_token) {
            console.error('❌ Login response missing required fields')
            return null
          }

          return {
            id: data.user.id,
            email: data.user.email,
            name: data.user.name,
            role: data.user.role,
            accessToken: data.access_token,
            refreshToken: data.refresh_token,
          }
        } catch (error) {
          console.error('❌ Error contacting auth service', error)
          throw new Error('Authentication service unavailable')
        }
      }
    })
  ],
  callbacks: {
    async jwt({ token, user }) {
      if (user) {
        token.accessToken = user.accessToken
        token.refreshToken = user.refreshToken
        token.role = user.role
      }
      return token
    },
    async session({ session, token }) {
      session.accessToken = token.accessToken
      session.refreshToken = token.refreshToken
      session.user.role = token.role
      return session
    }
  },
  pages: {
    signIn: '/login',
  },
  session: {
    strategy: 'jwt',
  },
}

const handler = NextAuth(authOptions)

export { handler as GET, handler as POST }
