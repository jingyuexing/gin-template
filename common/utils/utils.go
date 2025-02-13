package utils

// BoolToError converts a boolean condition to an error
// If the condition is true, returns nil
// If the condition is false, returns the provided error
func BoolToError(condition bool, err error) error {
    if condition {
        return nil
    }
    return err
}
