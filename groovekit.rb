class Groovekit < Formula
  desc "CLI for GrooveKit - Monitor cron jobs and APIs"
  homepage "https://groovekit.io"
  url "https://github.com/scookdev/groovekit-cli/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "PLACEHOLDER_SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", bin/"groovekit"
  end

  test do
    assert_match "groovekit version", shell_output("#{bin}/groovekit --version 2>&1", 1)
  end
end
