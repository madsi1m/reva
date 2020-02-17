<?php
// GENERATED CODE -- DO NOT EDIT!

namespace Revad\Filesender;

/**
 */
class FileSenderServiceClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \Revad\Filesender\HelloRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     */
    public function Hello(\Revad\Filesender\HelloRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/revad.filesender.FileSenderService/Hello',
        $argument,
        ['\Revad\Filesender\HelloResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Revad\Filesender\ReadChunkRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     */
    public function ReadChunk(\Revad\Filesender\ReadChunkRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/revad.filesender.FileSenderService/ReadChunk',
        $argument,
        ['\Revad\Filesender\ReadChunkResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Revad\Filesender\WriteChunkRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     */
    public function WriteChunk(\Revad\Filesender\WriteChunkRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/revad.filesender.FileSenderService/WriteChunk',
        $argument,
        ['\Revad\Filesender\WriteChunkResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Revad\Filesender\CompleteFileRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     */
    public function CompleteFile(\Revad\Filesender\CompleteFileRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/revad.filesender.FileSenderService/CompleteFile',
        $argument,
        ['\Revad\Filesender\CompleteFileResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Revad\Filesender\DeleteFileRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     */
    public function DeleteFile(\Revad\Filesender\DeleteFileRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/revad.filesender.FileSenderService/DeleteFile',
        $argument,
        ['\Revad\Filesender\DeleteFileResponse', 'decode'],
        $metadata, $options);
    }

    /**
     * @param \Revad\Filesender\StoreWholeFileRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     */
    public function StoreWholeFile(\Revad\Filesender\StoreWholeFileRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/revad.filesender.FileSenderService/StoreWholeFile',
        $argument,
        ['\Revad\Filesender\StoreWholeFileResponse', 'decode'],
        $metadata, $options);
    }

}
