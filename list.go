package redistructs

// func (rs RedigoStructs) List(ctx context.Context, dest types.RediStruct, mods ...rq.Modifier) error {
// 	dt := reflect.ValueOf(dest)
// 	if dt.Kind() != reflect.Ptr || dt.IsNil() {
// 		return errors.New("must pass a slice ptr")
// 	}
// 	dt = dt.Elem()
// 	if dt.Kind() != reflect.Slice {
// 		return errors.New("must pass a slice ptr")
// 	}

// 	conn, err := rs.pool.GetContext(ctx)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to acquire a connection")
// 	}
// 	defer conn.Close()

// 	dbIdx, changed := rs.getDBIndex(rs.model)
// 	if changed {
// 		conn.Do("SELECT", dbIdx)
// 	}

// 	keys, err := types.selectKeys(conn, rs.name, mods)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to select query")
// 	}

// 	if len(keys) == 0 {
// 		return nil
// 	}

// 	if s.HashStoreEnabled {
// 		for _, key := range keys {
// 			err := conn.Send("HGETALL", key)
// 			if err != nil {
// 				return errors.Wrapf(err, "failed to send HGETALL %s", key)
// 			}
// 		}
// 	} else {
// 		if len(s.model.Serialized()) == 0 {
// 			return errors.Errorf("failed to implement Serialized %v", dest)
// 		}

// 		args := []interface{}{s.KeyPrefix}
// 		for _, k := range keys {
// 			args = append(args, k)
// 		}
// 		err = conn.Send("HMGET", args...)
// 		if err != nil {
// 			return errors.Wrapf(err, "faild to send HMGET %s %s", s.KeyPrefix, keys)
// 		}
// 	}

// 	err = conn.Flush()
// 	if err != nil {
// 		return errors.Wrapf(err, "faild to FLUSH")
// 	}

// 	if s.HashStoreEnabled {
// 		vt := dt.Type().Elem().Elem()
// 		for _, key := range keys {
// 			v, err := redis.Values(conn.Receive())
// 			if err != nil {
// 				return errors.Wrap(err, "faild to receive or cast redis command result")
// 			}
// 			vv := reflect.New(vt)
// 			err = redis.ScanStruct(v, vv.Interface())
// 			if err != nil {
// 				return errors.Wrapf(err, "faild to scan struct %s %x", key, v)
// 			}
// 			dt.Set(reflect.Append(dt, vv))
// 		}
// 	} else {
// 		vt := dt.Type().Elem()
// 		v, err := redis.Values(conn.Receive())
// 		if err != nil {
// 			return errors.Wrap(err, fmt.Sprintf("faild to receive or cast redis command result. keys: %v", keys))
// 		}
// 		for _, w := range v {
// 			if w == nil {
// 				continue
// 			}
// 			vv := reflect.New(vt)
// 			vv.MethodByName("Deserialized").Call([]reflect.Value{reflect.ValueOf(w)})
// 			dt.Set(reflect.Append(dt, vv.Elem()))
// 		}
// 	}

// 	return nil
// }
