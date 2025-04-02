import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs.tsx';

export default function Settings() {
  return (
    <Tabs defaultValue="general">
      <TabsList>
        <TabsTrigger value="general">General</TabsTrigger>
        <TabsTrigger value="s3-options">S3 Options</TabsTrigger>
      </TabsList>
      <TabsContent value="general">
        <h2 className="text-xl font-bold">General Settings</h2>
        <p>General settings content goes here.</p>
      </TabsContent>
      <TabsContent value="s3-options">
        <h2 className="text-xl font-bold">S3 Options</h2>
        <p>S3 options content goes here.</p>
      </TabsContent>
    </Tabs>
  );
}