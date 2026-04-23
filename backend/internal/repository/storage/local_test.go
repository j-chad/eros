package storage

import (
	"backend/internal/testutil"
	"context"
	"io"
	"strings"
	"testing"
)

func TestNewLocalFileStore_EmptyRoot(t *testing.T) {
	_, err := NewLocalFileStore("")
	testutil.NotNilErr(t, err)
}

func TestNewLocalFileStore_CreatesDir(t *testing.T) {
	dir := t.TempDir() + "/subdir"
	store, err := NewLocalFileStore(dir)
	testutil.NilErr(t, err)
	testutil.NotNil(t, store)
}

func TestLocalFileStore_PutAndGet(t *testing.T) {
	store, err := NewLocalFileStore(t.TempDir())
	testutil.NilErr(t, err)

	ctx := context.Background()
	content := "hello world"
	key, err := store.Put(ctx, "test.txt", "text/plain", strings.NewReader(content))
	testutil.NilErr(t, err)
	testutil.True(t, len(key) > 0, "key should not be empty")

	reader, err := store.Get(ctx, key)
	testutil.NilErr(t, err)
	defer reader.Close()

	data, err := io.ReadAll(reader)
	testutil.NilErr(t, err)
	testutil.Equal(t, string(data), content)
}

func TestLocalFileStore_PutPreservesExtension(t *testing.T) {
	store, err := NewLocalFileStore(t.TempDir())
	testutil.NilErr(t, err)

	key, err := store.Put(context.Background(), "photo.jpg", "image/jpeg", strings.NewReader("img"))
	testutil.NilErr(t, err)
	testutil.True(t, strings.HasSuffix(key, ".jpg"), "key should preserve extension: "+key)
}

func TestLocalFileStore_Delete(t *testing.T) {
	store, err := NewLocalFileStore(t.TempDir())
	testutil.NilErr(t, err)

	ctx := context.Background()
	key, err := store.Put(ctx, "test.txt", "text/plain", strings.NewReader("data"))
	testutil.NilErr(t, err)

	testutil.NilErr(t, store.Delete(ctx, key))

	_, err = store.Get(ctx, key)
	testutil.NotNilErr(t, err)
}

func TestLocalFileStore_DeleteNonExistent(t *testing.T) {
	store, err := NewLocalFileStore(t.TempDir())
	testutil.NilErr(t, err)

	// Delete of non-existent file should not error (os.IsNotExist is ignored)
	err = store.Delete(context.Background(), "00000000-0000-4000-8000-000000000000.txt")
	testutil.NilErr(t, err)
}

func TestLocalFileStore_DeleteMany(t *testing.T) {
	store, err := NewLocalFileStore(t.TempDir())
	testutil.NilErr(t, err)

	ctx := context.Background()
	k1, _ := store.Put(ctx, "a.txt", "text/plain", strings.NewReader("a"))
	k2, _ := store.Put(ctx, "b.txt", "text/plain", strings.NewReader("b"))

	testutil.NilErr(t, store.DeleteMany(ctx, []string{k1, k2}))

	_, err = store.Get(ctx, k1)
	testutil.NotNilErr(t, err)
	_, err = store.Get(ctx, k2)
	testutil.NotNilErr(t, err)
}

func TestLocalFileStore_List(t *testing.T) {
	store, err := NewLocalFileStore(t.TempDir())
	testutil.NilErr(t, err)

	ctx := context.Background()
	store.Put(ctx, "a.txt", "text/plain", strings.NewReader("a"))
	store.Put(ctx, "b.txt", "text/plain", strings.NewReader("b"))

	var keys []string
	for key, err := range store.List(ctx) {
		testutil.NilErr(t, err)
		keys = append(keys, key)
	}
	testutil.Equal(t, len(keys), 2)
}

func TestSafePath_ValidKey(t *testing.T) {
	store := &LocalFileStore{root: "/tmp/files"}
	path, err := store.safePath("550e8400-e29b-41d4-a716-446655440000.jpg")
	testutil.NilErr(t, err)
	testutil.Equal(t, path, "/tmp/files/550e8400-e29b-41d4-a716-446655440000.jpg")
}

func TestSafePath_ValidKeyNoExtension(t *testing.T) {
	store := &LocalFileStore{root: "/tmp/files"}
	_, err := store.safePath("550e8400-e29b-41d4-a716-446655440000")
	testutil.NilErr(t, err)
}

func TestSafePath_RejectsTraversal(t *testing.T) {
	store := &LocalFileStore{root: "/tmp/files"}
	tests := []string{
		"../etc/passwd",
		"../../secret",
		"foo/bar",
		"",
		"not-a-uuid.txt",
		"550e8400-e29b-41d4-a716-446655440000.toolong",
	}
	for _, key := range tests {
		t.Run(key, func(t *testing.T) {
			_, err := store.safePath(key)
			testutil.NotNilErr(t, err)
		})
	}
}

func TestAllowedKeyRegex(t *testing.T) {
	tests := []struct {
		key  string
		want bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"550e8400-e29b-41d4-a716-446655440000.jpg", true},
		{"550e8400-e29b-41d4-a716-446655440000.mp4", true},
		{"550e8400-e29b-41d4-a716-446655440000.gz", true},
		{"550e8400-e29b-41d4-a716-446655440000.toolong", false},
		{"not-a-uuid", false},
		{"../traversal", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			testutil.Equal(t, ALLOWED_KEY_REGEX.MatchString(tt.key), tt.want)
		})
	}
}
