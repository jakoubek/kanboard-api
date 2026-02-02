package kanboard

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

// GetTimezone returns the server's configured timezone string (e.g., "UTC", "Europe/Berlin").
func (c *Client) GetTimezone(ctx context.Context) (string, error) {
	var tz string
	if err := c.call(ctx, "getTimezone", nil, &tz); err != nil {
		return "", fmt.Errorf("getTimezone: %w", err)
	}
	return tz, nil
}

// loadTimezone fetches and caches the timezone location from the server.
func (c *Client) loadTimezone(ctx context.Context) error {
	tz, err := c.GetTimezone(ctx)
	if err != nil {
		return err
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return fmt.Errorf("invalid timezone %q: %w", tz, err)
	}
	c.timezone = loc
	return nil
}

// ensureTimezone loads the timezone if tzEnabled and not yet loaded.
func (c *Client) ensureTimezone(ctx context.Context) error {
	if !c.tzEnabled {
		return nil
	}
	var err error
	c.tzOnce.Do(func() {
		err = c.loadTimezone(ctx)
	})
	return err
}

// convertTimestamps converts all Timestamp fields in v to the client's timezone.
// v must be a pointer. Handles structs, pointers to structs, and slices of structs.
func (c *Client) convertTimestamps(v any) {
	if c.timezone == nil {
		return
	}
	rv := reflect.ValueOf(v)
	c.walkAndConvert(rv)
}

var timestampType = reflect.TypeOf(Timestamp{})

func (c *Client) walkAndConvert(rv reflect.Value) {
	switch rv.Kind() {
	case reflect.Ptr:
		if !rv.IsNil() {
			c.walkAndConvert(rv.Elem())
		}
	case reflect.Struct:
		if rv.Type() == timestampType {
			if rv.CanSet() {
				ts := rv.Addr().Interface().(*Timestamp)
				if !ts.IsZero() {
					ts.Time = ts.Time.In(c.timezone)
				}
			}
			return
		}
		for i := 0; i < rv.NumField(); i++ {
			c.walkAndConvert(rv.Field(i))
		}
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			c.walkAndConvert(rv.Index(i))
		}
	}
}
