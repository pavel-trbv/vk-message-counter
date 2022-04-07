import { FC } from 'react';
import { Heading, Link, Text } from '@chakra-ui/react';

const About: FC = () => {
  return (
    <>
      <Heading mb="4">About</Heading>
      <Text lineHeight="6" mb="4">
        This is a simple application that allows you to get statistics of chat messages on the&nbsp;
        <Link href="https://vk.com" color="blue.500">
          vk.com
        </Link>{' '}
        social network.
        <br />
        <br />
        It was created using the Next.js framework and Chakra UI library.
        <br />
        <br />
        The server side of this application is written in Golang in the form of HTTP API. Counting the number of messages is carried out using the VK
        API.
        <br />
        <br />
        Source code:{' '}
        <Link href="https://github.com/pavel-trbv/https://github.com/pavel-trbv/vk-message-counter" color="blue.500">
          vk-message-counter{' '}
        </Link>
      </Text>
    </>
  );
};

export default About;
