import { Input } from '@/components/ui/input.tsx';
import { Label } from '@/components/ui/label.tsx';
import { useEffect, useState } from 'react';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select.tsx';
import { Button } from '@/components/ui/button.tsx';
import { toast } from 'sonner';
import React from 'react';

export default function Home() {
  const [path, setPath] = useState<string>('');
  const [bucket, setBucket] = useState<string>('');
  const [bucketPath, setBucketPath] = useState<string>('');
  const [buckets, setBuckets] = useState<string[]>([]);

  useEffect(() => {
    fetch('/api/buckets')
      .then((res) => res.json())
      .then((data) => {
        if (Array.isArray(data)) {
          setBuckets(data);
        } else if (data.error) {
          toast(data.error);
        }
      })
      .catch((error) => {
        console.error('Error fetching buckets:', error);
        toast(error);
      });
  }, []);
  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    fetch('/api/backup', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        path,
        bucket,
        bucket_path: bucketPath,
      }),
    })
      .then(async (res) => {
        if (!res.ok) {
          const data = await res.json();
          if (data.error)
            return toast(data.error);
        }
      })
      .catch((error) => {
        console.error('Error starting backup:', error);
        toast(error);
      });
  };

  return (
    <div>
      <form onSubmit={onSubmit} className="space-y-4">
        <Label htmlFor="path">What do you want to backup?</Label>
        <Input placeholder="/opt/data..." id="path" value={path} onChange={(e) => setPath(e.target.value)} />
        <Label htmlFor="path">Where do you want to backup?</Label>
        <div className="flex gap-2">
          <Select onValueChange={setBucket} defaultValue={bucket}>
            <SelectTrigger>
              <SelectValue placeholder="Select a bucket" />
            </SelectTrigger>
            <SelectContent>
              {buckets.map((bucket) => (
                <SelectItem key={bucket} value={bucket}>
                  {bucket}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Input placeholder="my/data" id="path" value={bucketPath} onChange={(e) => setBucketPath(e.target.value)} />
        </div>
        <Button type="submit">
          Start Backup
        </Button>
      </form>
    </div>
  );
}