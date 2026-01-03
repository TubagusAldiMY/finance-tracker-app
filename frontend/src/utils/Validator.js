export class Validator {
  constructor(props, rules) {
    this.props = props        // data form
    this.rules = rules        // aturan validasi
    this.errors = {}          // hasil validasi
  }

  validate() {
    this.errors = {}

    for (const field in this.rules) {
      const value = this.props[field]
      const fieldRules = this.rules[field]

      for (const rule of fieldRules) {
        const error = rule(value, this.props)

        if (error) {
          this.errors[field] = error
          break // stop di rule pertama (best practice UX)
        }
      }
    }

    return Object.keys(this.errors).length === 0
  }

  getErrors() {
    return this.errors
  }

  hasError(field) {
    return !!this.errors[field]
  }

  getError(field) {
    return this.errors[field] || ''
  }
}
export const required = (message = 'Required') => value =>
  value ? '' : message

export const email = (message = 'Invalid email') => value =>
  /^\S+@\S+\.\S+$/.test(value) ? '' : message

export const minLength = (min, message) => value =>
  value?.length >= min ? '' : message || `Minimum ${min} characters`

export const sameAs = (field, message) => (value, props) =>
  value === props[field] ? '' : message
